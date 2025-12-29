package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Endpoint   string
	Model      string
	httpClient *http.Client
}

func NewClient(endpoint, model string) *Client {
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	if model == "" {
		model = "mshop"
	}
	return &Client{
		Endpoint: endpoint,
		Model:    model,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func (c *Client) Generate(context context.Context, prompt string, maxTokens int) (string, error) {
	payload := map[string]interface{}{
		"model":       c.Model,
		"prompt":      prompt,
		"stream":      false,
		"temperature": 0.1,
	}
	if maxTokens > 0 {
		payload["max_tokens"] = maxTokens
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context, "POST", c.Endpoint+"/api/generate", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to ollama failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var res struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return string(body), nil
	}

	stringBody := res.Response
	start := strings.IndexByte(stringBody, '{')
	end := strings.LastIndexByte(stringBody, '}')

	if start == -1 || end == -1 || end < start {
		return `{"intent":"UNKNOWN"}`, nil
	}

	return strings.TrimSpace(res.Response), nil
}

package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mshop/ai-service/ai"
	token "github.com/mshop/ai-service/auth"
	"github.com/mshop/ai-service/models"
)

type AIHandler struct {
	client *ai.Client
}

func NewAIHandler(client *ai.Client) *AIHandler {
	return &AIHandler{client: client}
}

func (h *AIHandler) Chat(c *gin.Context) {

	if os.Getenv("DISABLE_AUTH") == "true" {
	} else {
		_, err := token.GetTokenClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}

	var prompt string
	var maxTokens int

	if p := c.Query("prompt"); p != "" {
		prompt = p
		if mt := c.Query("max_tokens"); mt != "" {
			if v, err := strconv.Atoi(mt); err == nil {
				maxTokens = v
			}
		}
	} else {
		var req models.AIRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing prompt (query or json) or invalid json"})
			return
		}
		prompt = req.Prompt
		maxTokens = req.MaxTokens
	}

	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt cannot be empty"})
		return
	}

	context, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	responseText, err := h.client.Generate(context, prompt, maxTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AIResponse{
		Response: responseText,
	})
}

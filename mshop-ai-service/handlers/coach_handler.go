package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mshop/ai-service/ai"
	token "github.com/mshop/ai-service/auth"
	"github.com/mshop/ai-service/models"
)

type CoachHandler struct {
	client *ai.Client
}

func NewCoachHandler(client *ai.Client) *CoachHandler {
	return &CoachHandler{client: client}
}

func (h *CoachHandler) Chat(c *gin.Context) {

	if os.Getenv("DISABLE_AUTH") != "true" {
		if _, err := token.GetTokenClaims(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}

	var req models.AIRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing prompt"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	coachPrompt := ai.BuildCoachPrompt(req.Prompt)

	raw, err := h.client.Generate(ctx, coachPrompt, req.MaxTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cleanJSON := ai.ExtractAndValidateCoachJSON(raw)

	c.JSON(http.StatusOK, models.AIResponse{
		Response: cleanJSON,
	})
}

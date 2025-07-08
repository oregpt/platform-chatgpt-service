package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/config"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/handlers"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/openai"
	"github.com/sirupsen/logrus"
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, openaiClient openai.ClientInterface, log *logrus.Logger, cfg *config.Config) {
	// Create handler
	handler := handlers.NewChatHandler(openaiClient, log, cfg)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "chatgpt-service",
		})
	})

	// API endpoints
	api := router.Group("/api")
	{
		// Chat endpoint
		api.POST("/chat", handler.HandleChat)
	}
}

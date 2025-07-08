package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/config"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/models"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/openai"
	"github.com/sirupsen/logrus"
)

// ChatHandler handles chat requests
type ChatHandler struct {
	openaiClient openai.ClientInterface
	log          *logrus.Logger
	cfg          *config.Config
}

// NewChatHandler creates a new chat handler
func NewChatHandler(openaiClient openai.ClientInterface, log *logrus.Logger, cfg *config.Config) *ChatHandler {
	return &ChatHandler{
		openaiClient: openaiClient,
		log:          log,
		cfg:          cfg,
	}
}

// HandleChat handles chat requests
func (h *ChatHandler) HandleChat(c *gin.Context) {
	startTime := time.Now()

	// Parse request
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Status: "error",
			Error: &models.ErrorInfo{
				Code:    "invalid_request",
				Message: "Invalid request format",
				Details: err.Error(),
			},
		})
		return
	}

	// Validate request
	if err := h.validateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Status: "error",
			Error: &models.ErrorInfo{
				Code:    "validation_error",
				Message: "Request validation failed",
				Details: err.Error(),
			},
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.cfg.RequestTimeout)
	defer cancel()

	// Process the chat request
	response, err := h.processChat(ctx, &req)
	if err != nil {
		h.log.Errorf("Error processing chat: %v", err)
		c.JSON(http.StatusInternalServerError, models.ChatResponse{
			Status:    "error",
			SessionID: req.SessionID,
			Error: &models.ErrorInfo{
				Code:    "processing_error",
				Message: "Error processing chat request",
				Details: err.Error(),
			},
		})
		return
	}

	// Calculate processing time
	processingTime := time.Since(startTime).Seconds()
	response.Metadata.ProcessingTime = processingTime
	response.Metadata.RequestID = req.Metadata.RequestID

	c.JSON(http.StatusOK, response)
}

// validateRequest validates the chat request
func (h *ChatHandler) validateRequest(req *models.ChatRequest) error {
	if req.OrganizationID == "" {
		return fmt.Errorf("organizationId is required")
	}
	if req.AgentID == "" {
		return fmt.Errorf("agentId is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("userId is required")
	}
	if req.Message == "" {
		return fmt.Errorf("message is required")
	}
	if req.SessionID == "" {
		return fmt.Errorf("sessionId is required")
	}
	if req.Context.AgentConfig.AIProvider != "chatgpt" {
		return fmt.Errorf("aiProvider must be 'chatgpt'")
	}
	return nil
}

// processChat processes a chat request
func (h *ChatHandler) processChat(ctx context.Context, req *models.ChatRequest) (*models.ChatResponse, error) {
	// Get or create thread (conversation)
	thread, err := h.openaiClient.GetOrCreateThread(ctx, req.SessionID, req.AgentID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create thread: %w", err)
	}

	// Add message to thread
	if err := h.openaiClient.AddMessageToThread(ctx, thread.ThreadID, req.Message); err != nil {
		return nil, fmt.Errorf("failed to add message to thread: %w", err)
	}

	// Run the thread with the specified model (using default model from config)
	response, err := h.openaiClient.RunThread(ctx, thread.ThreadID, h.cfg.DefaultModel)
	if err != nil {
		return nil, fmt.Errorf("failed to run thread: %w", err)
	}

	// Create response
	chatResponse := &models.ChatResponse{
		Response:       response,
		SessionID:      req.SessionID,
		ConversationID: thread.ThreadID, // Use thread ID as conversation ID
		Status:         "success",
		Metadata: models.ResponseMeta{
			Model:      h.cfg.DefaultModel,
			TokensUsed: 0, // We don't track token usage in this implementation
			Provider:   "chatgpt",
			Cost:       0, // We don't track cost in this implementation
		},
		Context: &models.ResponseContext{
			ThreadID:    thread.ThreadID,
			AssistantID: "", // No assistant ID in this implementation
		},
	}

	return chatResponse, nil
}

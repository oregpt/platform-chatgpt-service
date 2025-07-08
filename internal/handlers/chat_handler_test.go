package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/config"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/models"
	"github.com/oregpt/agentplatform-chatgpt-service/internal/openai"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	// Setup
	log := logrus.New()
	cfg := &config.Config{}
	openaiClient := openai.NewMockClient(log)
	handler := NewChatHandler(openaiClient, log, cfg)

	// Test cases
	testCases := []struct {
		name        string
		request     models.ChatRequest
		expectError bool
	}{
		{
			name: "Valid request",
			request: models.ChatRequest{
				OrganizationID: "org123",
				AgentID:        "agent123",
				UserID:         "user123",
				Message:        "Hello",
				SessionID:      "session123",
				Context: models.Context{
					AgentConfig: models.AgentConfig{
						AIProvider: "chatgpt",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Missing organization ID",
			request: models.ChatRequest{
				AgentID:   "agent123",
				UserID:    "user123",
				Message:   "Hello",
				SessionID: "session123",
				Context: models.Context{
					AgentConfig: models.AgentConfig{
						AIProvider: "chatgpt",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Missing agent ID",
			request: models.ChatRequest{
				OrganizationID: "org123",
				UserID:         "user123",
				Message:        "Hello",
				SessionID:      "session123",
				Context: models.Context{
					AgentConfig: models.AgentConfig{
						AIProvider: "chatgpt",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Wrong AI provider",
			request: models.ChatRequest{
				OrganizationID: "org123",
				AgentID:        "agent123",
				UserID:         "user123",
				Message:        "Hello",
				SessionID:      "session123",
				Context: models.Context{
					AgentConfig: models.AgentConfig{
						AIProvider: "claude",
					},
				},
			},
			expectError: true,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := handler.validateRequest(&tc.request)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleChatInvalidRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	log := logrus.New()
	cfg := &config.Config{}
	openaiClient := openai.NewMockClient(log)
	handler := NewChatHandler(openaiClient, log, cfg)

	// Create router
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Create invalid request
	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "invalid_request", response.Error.Code)
}

func TestHandleChatValidationError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	log := logrus.New()
	cfg := &config.Config{}
	openaiClient := openai.NewMockClient(log)
	handler := NewChatHandler(openaiClient, log, cfg)

	// Create router
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Create request with missing required fields
	chatRequest := models.ChatRequest{
		Message: "Hello",
	}
	requestBody, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "validation_error", response.Error.Code)
}

func TestHandleChatSuccess(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	log := logrus.New()
	cfg := &config.Config{
		RequestTimeout: 30 * time.Second,
	}
	
	// Create mock client with custom behavior
	mockClient := openai.NewMockClient(log)
	mockClient.RunThreadFunc = func(ctx context.Context, threadID, model string) (string, error) {
		return "This is a test response", nil
	}
	
	handler := NewChatHandler(mockClient, log, cfg)

	// Create router
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Create valid request
	chatRequest := models.ChatRequest{
		OrganizationID: "org123",
		AgentID:        "agent123",
		UserID:         "user123",
		Message:        "Hello",
		SessionID:      "session123",
		Context: models.Context{
			AgentConfig: models.AgentConfig{
				AIProvider: "chatgpt",
				Name:       "Test Assistant",
			},
		},
		Metadata: models.Metadata{
			RequestID: "req123",
		},
	}
	requestBody, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "This is a test response", response.Response)
	assert.Equal(t, "session123", response.SessionID)
	assert.Equal(t, "mock-thread-id", response.ConversationID)
	assert.Equal(t, "chatgpt", response.Metadata.Provider)
	assert.Equal(t, "req123", response.Metadata.RequestID)
}

func TestHandleChatError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	log := logrus.New()
	cfg := &config.Config{
		RequestTimeout: 30 * time.Second,
	}
	
	// Create mock client with error behavior
	mockClient := openai.NewMockClient(log)
	mockClient.RunThreadFunc = func(ctx context.Context, threadID, model string) (string, error) {
		return "", errors.New("API error")
	}
	
	handler := NewChatHandler(mockClient, log, cfg)

	// Create router
	router := gin.New()
	router.POST("/chat", handler.HandleChat)

	// Create valid request
	chatRequest := models.ChatRequest{
		OrganizationID: "org123",
		AgentID:        "agent123",
		UserID:         "user123",
		Message:        "Hello",
		SessionID:      "session123",
		Context: models.Context{
			AgentConfig: models.AgentConfig{
				AIProvider: "chatgpt",
			},
		},
	}
	requestBody, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ChatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "processing_error", response.Error.Code)
	assert.Contains(t, response.Error.Details, "API error")
}

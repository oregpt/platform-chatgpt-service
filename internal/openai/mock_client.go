package openai

import (
	"context"
	"time"

	"github.com/oregpt/agentplatform-chatgpt-service/internal/models"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

// MockClient is a mock implementation of the OpenAI client for testing
type MockClient struct {
	GetOrCreateThreadFunc func(ctx context.Context, sessionID, agentID, userID string) (*models.ThreadInfo, error)
	AddMessageToThreadFunc func(ctx context.Context, threadID, content string) error
	RunThreadFunc func(ctx context.Context, threadID, model string) (string, error)
	CleanupOldCacheEntriesFunc func(threadTTL time.Duration)
}

// NewMockClient creates a new mock OpenAI client
func NewMockClient(log *logrus.Logger) *MockClient {
	return &MockClient{
		GetOrCreateThreadFunc: func(ctx context.Context, sessionID, agentID, userID string) (*models.ThreadInfo, error) {
			return &models.ThreadInfo{
				ThreadID:  "mock-thread-id",
				SessionID: sessionID,
				AgentID:   agentID,
				UserID:    userID,
				Messages:  []openai.ChatCompletionMessage{},
				CreatedAt: time.Now(),
				LastUsed:  time.Now(),
			}, nil
		},
		AddMessageToThreadFunc: func(ctx context.Context, threadID, content string) error {
			return nil
		},
		RunThreadFunc: func(ctx context.Context, threadID, model string) (string, error) {
			return "This is a mock response from the OpenAI API.", nil
		},
		CleanupOldCacheEntriesFunc: func(threadTTL time.Duration) {
			// Do nothing in mock
		},
	}
}

// GetOrCreateThread gets an existing thread or creates a new one
func (c *MockClient) GetOrCreateThread(ctx context.Context, sessionID, agentID, userID string) (*models.ThreadInfo, error) {
	return c.GetOrCreateThreadFunc(ctx, sessionID, agentID, userID)
}

// AddMessageToThread adds a message to a thread
func (c *MockClient) AddMessageToThread(ctx context.Context, threadID, content string) error {
	return c.AddMessageToThreadFunc(ctx, threadID, content)
}

// RunThread runs a thread with the model and returns the assistant's response
func (c *MockClient) RunThread(ctx context.Context, threadID, model string) (string, error) {
	return c.RunThreadFunc(ctx, threadID, model)
}

// CleanupOldCacheEntries removes old entries from the cache
func (c *MockClient) CleanupOldCacheEntries(threadTTL time.Duration) {
	c.CleanupOldCacheEntriesFunc(threadTTL)
}

package openai

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/oregpt/agentplatform-chatgpt-service/internal/models"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

// Client wraps the OpenAI client with additional functionality
type Client struct {
	client      *openai.Client
	log         *logrus.Logger
	threadCache map[string]*models.ThreadInfo
	threadMutex sync.RWMutex
}

// NewClient creates a new OpenAI client wrapper
func NewClient(apiKey string, log *logrus.Logger) *Client {
	return &Client{
		client:      openai.NewClient(apiKey),
		log:         log,
		threadCache: make(map[string]*models.ThreadInfo),
	}
}

// GetOrCreateThread gets an existing thread or creates a new one
func (c *Client) GetOrCreateThread(ctx context.Context, sessionID, agentID, userID string) (*models.ThreadInfo, error) {
	// Check cache first
	c.threadMutex.RLock()
	thread, exists := c.threadCache[sessionID]
	c.threadMutex.RUnlock()

	if exists {
		// Update last used time
		c.threadMutex.Lock()
		thread.LastUsed = time.Now()
		c.threadMutex.Unlock()
		return thread, nil
	}

	// Create a new thread
	c.log.Infof("Creating new thread for session %s", sessionID)
	
	// Create the thread info
	threadInfo := &models.ThreadInfo{
		ThreadID:  sessionID, // Use sessionID as threadID for simplicity
		SessionID: sessionID,
		AgentID:   agentID,
		UserID:    userID,
		Messages:  []openai.ChatCompletionMessage{},
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	
	c.threadMutex.Lock()
	c.threadCache[sessionID] = threadInfo
	c.threadMutex.Unlock()
	
	return threadInfo, nil
}

// AddMessageToThread adds a message to a thread
func (c *Client) AddMessageToThread(ctx context.Context, threadID, content string) error {
	c.threadMutex.Lock()
	defer c.threadMutex.Unlock()
	
	thread, exists := c.threadCache[threadID]
	if !exists {
		return fmt.Errorf("thread %s not found", threadID)
	}
	
	// Add user message to the thread
	thread.Messages = append(thread.Messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: content,
	})
	
	return nil
}

// RunThread runs a thread with the model and returns the assistant's response
func (c *Client) RunThread(ctx context.Context, threadID, model string) (string, error) {
	c.threadMutex.RLock()
	thread, exists := c.threadCache[threadID]
	if !exists {
		c.threadMutex.RUnlock()
		return "", fmt.Errorf("thread %s not found", threadID)
	}
	
	// Make a copy of the messages to avoid race conditions
	messages := make([]openai.ChatCompletionMessage, len(thread.Messages))
	copy(messages, thread.Messages)
	c.threadMutex.RUnlock()
	
	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}
	
	// Call the OpenAI API
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices returned")
	}
	
	// Get the assistant's response
	assistantResponse := resp.Choices[0].Message.Content
	
	// Add the assistant's response to the thread
	c.threadMutex.Lock()
	thread, stillExists := c.threadCache[threadID]
	if stillExists {
		thread.Messages = append(thread.Messages, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: assistantResponse,
		})
	}
	c.threadMutex.Unlock()
	
	return assistantResponse, nil
}

// CleanupOldCacheEntries removes old entries from the cache
func (c *Client) CleanupOldCacheEntries(threadTTL time.Duration) {
	now := time.Now()
	
	// Cleanup threads
	c.threadMutex.Lock()
	for sessionID, thread := range c.threadCache {
		if now.Sub(thread.LastUsed) > threadTTL {
			c.log.Infof("Removing thread %s from cache due to TTL expiration", thread.ThreadID)
			delete(c.threadCache, sessionID)
		}
	}
	c.threadMutex.Unlock()
}

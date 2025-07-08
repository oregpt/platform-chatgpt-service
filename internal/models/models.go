package models

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

// ChatRequest represents the standardized input schema for all AI services
type ChatRequest struct {
	OrganizationID string    `json:"organizationId"`
	AgentID        string    `json:"agentId"`
	UserID         string    `json:"userId"`
	Message        string    `json:"message"`
	SessionID      string    `json:"sessionId"`
	Context        Context   `json:"context"`
	Metadata       Metadata  `json:"metadata"`
}

// Context represents the context information for the chat request
type Context struct {
	Files       []File       `json:"files"`
	ChatHistory []ChatEntry  `json:"chatHistory"`
	AgentConfig AgentConfig  `json:"agentConfig"`
}

// File represents a file in the context
type File struct {
	Filename     string    `json:"filename"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
}

// ChatEntry represents a single entry in the chat history
type ChatEntry struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// AgentConfig represents the configuration for an agent
type AgentConfig struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Instructions string  `json:"instructions"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"maxTokens"`
	AIProvider   string  `json:"aiProvider"` // Should be "chatgpt" for this service
}

// Metadata represents metadata for the request
type Metadata struct {
	RequestID  string    `json:"requestId"`
	Timestamp  time.Time `json:"timestamp"`
	UserAgent  string    `json:"userAgent"`
}

// ChatResponse represents the standardized output schema for all AI services
type ChatResponse struct {
	Response       string         `json:"response"`
	SessionID      string         `json:"sessionId"`
	ConversationID string         `json:"conversationId"`
	Status         string         `json:"status"` // "success", "error", or "timeout"
	Metadata       ResponseMeta   `json:"metadata"`
	Error          *ErrorInfo     `json:"error,omitempty"`
	Context        *ResponseContext `json:"context,omitempty"`
}

// ResponseMeta represents metadata for the response
type ResponseMeta struct {
	Model          string  `json:"model"`
	TokensUsed     int     `json:"tokensUsed"`
	ProcessingTime float64 `json:"processingTime"`
	Provider       string  `json:"provider"` // Should be "chatgpt" for this service
	Cost           float64 `json:"cost"`
	RequestID      string  `json:"requestId"`
}

// ErrorInfo represents error information in the response
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ResponseContext represents additional context information in the response
type ResponseContext struct {
	ThreadID     string   `json:"threadId"`
	AssistantID  string   `json:"assistantId"`
	NextActions  []string `json:"nextActions,omitempty"`
}

// AssistantInfo represents information about an OpenAI Assistant
type AssistantInfo struct {
	AssistantID  string
	AgentID      string
	Instructions string
	Model        string
	FileIDs      []string
	CreatedAt    time.Time
	LastUsed     time.Time
}

// ThreadInfo represents information about a chat thread
type ThreadInfo struct {
	ThreadID     string
	SessionID    string
	AgentID      string
	UserID       string
	Messages     []openai.ChatCompletionMessage
	CreatedAt    time.Time
	LastUsed     time.Time
}

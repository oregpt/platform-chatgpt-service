package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	OpenAIAPIKey   string
	Port           string
	ThreadTTL      time.Duration
	DefaultModel   string
	MaxRetries     int
	RetryDelay     time.Duration
	RequestTimeout time.Duration
}

// NewConfig creates a new configuration with values from environment variables
func NewConfig() *Config {
	// Get OpenAI API key from environment
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}



	// Get thread TTL from environment or use default (60 minutes)
	threadTTLStr := os.Getenv("THREAD_TTL")
	threadTTL := 60 * time.Minute
	if threadTTLStr != "" {
		if ttl, err := strconv.Atoi(threadTTLStr); err == nil {
			threadTTL = time.Duration(ttl) * time.Minute
		}
	}

	// Get default model from environment or use default
	defaultModel := os.Getenv("DEFAULT_MODEL")
	if defaultModel == "" {
		defaultModel = "gpt-4o"
	}

	// Get max retries from environment or use default
	maxRetriesStr := os.Getenv("MAX_RETRIES")
	maxRetries := 3
	if maxRetriesStr != "" {
		if mr, err := strconv.Atoi(maxRetriesStr); err == nil {
			maxRetries = mr
		}
	}

	// Get retry delay from environment or use default (1 second)
	retryDelayStr := os.Getenv("RETRY_DELAY")
	retryDelay := 1 * time.Second
	if retryDelayStr != "" {
		if rd, err := strconv.Atoi(retryDelayStr); err == nil {
			retryDelay = time.Duration(rd) * time.Second
		}
	}

	// Get request timeout from environment or use default (60 seconds)
	requestTimeoutStr := os.Getenv("REQUEST_TIMEOUT")
	requestTimeout := 60 * time.Second
	if requestTimeoutStr != "" {
		if rt, err := strconv.Atoi(requestTimeoutStr); err == nil {
			requestTimeout = time.Duration(rt) * time.Second
		}
	}

	return &Config{
		OpenAIAPIKey:   openAIAPIKey,
		Port:           port,
		ThreadTTL:      threadTTL,
		DefaultModel:   defaultModel,
		MaxRetries:     maxRetries,
		RetryDelay:     retryDelay,
		RequestTimeout: requestTimeout,
	}
}

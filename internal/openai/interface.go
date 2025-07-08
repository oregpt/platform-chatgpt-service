package openai

import (
	"context"
	"time"

	"github.com/oregpt/agentplatform-chatgpt-service/internal/models"
)

// ClientInterface defines the interface for the OpenAI client
type ClientInterface interface {
	GetOrCreateThread(ctx context.Context, sessionID, agentID, userID string) (*models.ThreadInfo, error)
	AddMessageToThread(ctx context.Context, threadID, content string) error
	RunThread(ctx context.Context, threadID, model string) (string, error)
	CleanupOldCacheEntries(threadTTL time.Duration)
}

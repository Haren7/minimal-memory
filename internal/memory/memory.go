package memory

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	Store(ctx context.Context, conversationID string, memories []Memory) error
	Retrieve(ctx context.Context, conversationID string, lastK int) ([]Memory, error)
}

type SemanticServiceInterface interface {
	Store(ctx context.Context, convesationID uuid.UUID, memories []Memory) (string, error)
	Retrieve(ctx context.Context, conversationID uuid.UUID, query string, topK int) ([]Memory, error)
	RetrieveSimilar(ctx context.Context, conversationID uuid.UUID, query string, topK int, rerankerOpts ...RerankerOpts) ([]Memory, error)
	RetrieveEpisode(ctx context.Context, conversationID uuid.UUID, episode time.Time, query string, topK int) ([]Memory, error)
}

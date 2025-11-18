package persistence

import (
	"context"
)

type RdbmsConversationRepoInterface interface {
	FetchOne(ctx context.Context, agent, user string) (RdbmsConversation, error)
	InsertOne(ctx context.Context, agent, user string) (RdbmsConversation, error)
}

type RdbmsMemoryRepoInterface interface {
	FetchOne(ctx context.Context, conversationID string) (RdbmsMemory, error)
	FetchMany(ctx context.Context, memoryIds []int) ([]RdbmsMemory, error)
	InsertOne(ctx context.Context, conversationID string, query, response string) (RdbmsMemory, error)
}

type VectorMemoryRepoInterface interface {
	Index(ctx context.Context, conversationID string, query, response string) error
	Search(ctx context.Context, conversationID string, query string, topK int) ([]VectorMemory, error)
}

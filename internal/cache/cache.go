package cache

import (
	"context"
)

type MemoryRepoInterface interface {
	SetOne(ctx context.Context, conversationID string, query string, response string) error
	SetMany(ctx context.Context, conversationID string, memories []Memory) error
	Get(ctx context.Context, conversationID string, lastK int) ([]Memory, error)
	Len(ctx context.Context, convesationID string) (int, error)
}

type MemoryRepo struct {
	memories map[string][]Memory
}

func NewMemoryRepo(size int) MemoryRepoInterface {
	return &MemoryRepo{
		memories: make(map[string][]Memory),
	}
}

func (r *MemoryRepo) SetOne(ctx context.Context, conversationID string, query string, response string) error {
	// if len > size, remove oldest
	return nil
}

func (r *MemoryRepo) SetMany(ctx context.Context, conversationID string, memories []Memory) error {
	// if len > size, remove oldest
	return nil
}

func (r *MemoryRepo) Get(ctx context.Context, conversationID string, lastK int) ([]Memory, error) {
	return nil, nil
}

func (r *MemoryRepo) Len(ctx context.Context, convesationID string) (int, error) {
	return 0, nil
}

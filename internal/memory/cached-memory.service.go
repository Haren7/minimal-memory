package memory

import (
	"context"
	"memory/internal/cache"
)

type CachedService struct {
	memoryRepo cache.MemoryRepoInterface
}

func NewCachedService(memoryRepo cache.MemoryRepoInterface) ServiceInterface {
	return &CachedService{
		memoryRepo: memoryRepo,
	}
}

func (r *CachedService) Store(ctx context.Context, conversationID string, memories []Memory) error {
	var cachedMemories []cache.Memory
	for _, memory := range memories {
		cachedMemories = append(cachedMemories, cache.Memory{
			Query:     memory.Query,
			Response:  memory.Response,
			CreatedAt: memory.CreatedAt,
		})
	}
	err := r.memoryRepo.SetMany(ctx, conversationID, cachedMemories)
	if err != nil {
		return err
	}
	return nil
}

func (r *CachedService) Retrieve(ctx context.Context, conversationID string, lastK int) ([]Memory, error) {
	cachedMemories, err := r.memoryRepo.Get(ctx, conversationID, lastK)

	if err != nil {
		return nil, err
	}

	var memories []Memory
	for _, cachedMemory := range cachedMemories {
		memories = append(memories, Memory{
			Query:     cachedMemory.Query,
			Response:  cachedMemory.Response,
			CreatedAt: cachedMemory.CreatedAt,
		})
	}

	return memories, nil
}

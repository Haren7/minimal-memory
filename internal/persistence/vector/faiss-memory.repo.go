package vector

import (
	"context"
	"memory/internal/embedding"
	"memory/internal/persistence"
)

/*
IDEA - folder per conversationID, for each op read the folder, index, write it back
If the conversation is not present, then create its files

what does
*/

type FaissMemoryRepo struct {
	faissClient     *FaissClient
	embeddingClient embedding.Service
	rdbmsMemoryRepo persistence.RdbmsMemoryRepoInterface
}

func NewFaissMemoryRepo(faissClient *FaissClient) persistence.VectorMemoryRepoInterface {
	return &FaissMemoryRepo{
		faissClient: faissClient,
	}
}

func (r *FaissMemoryRepo) Index(ctx context.Context, conversationID string, query, response string) error {
	memory, err := r.rdbmsMemoryRepo.InsertOne(ctx, conversationID, query, response)
	if err != nil {
		return err
	}
	memoryId := memory.ID
	embedding, err := r.embeddingClient.EmbedOne(ctx, query)
	if err != nil {
		return err
	}
	err = r.faissClient.Index(ctx, conversationID, memoryId, embedding)
	if err != nil {
		return err
	}
	return nil
}

func (r *FaissMemoryRepo) Search(ctx context.Context, conversationID string, query string, topK int) ([]persistence.VectorMemory, error) {
	embedding, err := r.embeddingClient.EmbedOne(ctx, query)
	if err != nil {
		return nil, err
	}
	faissResponse, err := r.faissClient.Search(ctx, conversationID, embedding, topK)
	var memoryIds []int
	for _, id := range faissResponse.Labels {
		memoryIds = append(memoryIds, int(id))
	}
	rdbmsMemories, err := r.rdbmsMemoryRepo.FetchMany(ctx, memoryIds)
	if err != nil {
		return nil, err
	}
	var vectorMemories []persistence.VectorMemory
	for _, memory := range rdbmsMemories {
		vectorMemories = append(vectorMemories, persistence.VectorMemory{
			UUID:      memory.UUID,
			Query:     memory.Query,
			Response:  memory.Response,
			CreatedAt: memory.CreatedAt,
		})
	}
	return vectorMemories, nil
}

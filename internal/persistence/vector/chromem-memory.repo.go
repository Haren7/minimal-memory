package vector

import (
	"context"
	"memory/internal/persistence"

	"github.com/philippgille/chromem-go"
)

type ChromemMemoryRepo struct {
	db *chromem.DB
}

func NewChromemMemoryRepo(db *chromem.DB) persistence.VectorMemoryRepoInterface {
	return &ChromemMemoryRepo{
		db: db,
	}
}

func (r *ChromemMemoryRepo) Index(ctx context.Context, conversationID, query, response string) error {
	collection, err := r.db.GetOrCreateCollection(conversationID, nil, nil)
	// collection.
	return nil
}

func (r *ChromemMemoryRepo) Search(ctx context.Context, conversationID, query string, topK int) ([]persistence.VectorMemory, error) {
	return nil, nil
}

package rdbms

import (
	"context"
	"database/sql"
	"memory/internal/persistence"
)

type MemoryRepo struct {
	db *sql.DB
}

func NewMemoryRepo(db *sql.DB) persistence.RdbmsMemoryRepoInterface {
	return &MemoryRepo{db}
}

func (r *MemoryRepo) FetchOne(ctx context.Context, conversationID string) (persistence.RdbmsMemory, error) {
	return persistence.RdbmsMemory{}, nil
}

func (r *MemoryRepo) FetchMany(ctx context.Context, memoryIds []int) ([]persistence.RdbmsMemory, error) {
	return nil, nil
}

func (r *MemoryRepo) InsertOne(ctx context.Context, conversationID string, query string, response string) (persistence.RdbmsMemory, error) {
	return persistence.RdbmsMemory{}, nil
}

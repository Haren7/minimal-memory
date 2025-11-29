package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/haren7/minimal-memory/internal/persistence"

	"github.com/google/uuid"
)

type MemoryRepo struct {
	db        *sql.DB
	tableName string
}

func NewMemoryRepo(db *sql.DB) persistence.MemoryRepoInterface {
	return &MemoryRepo{db: db, tableName: "memories"}
}

func NewFaissMemoryRepo(db *sql.DB) persistence.MemoryRepoInterface {
	return &MemoryRepo{db: db, tableName: "memories_meta"}
}

func (r *MemoryRepo) FetchOne(ctx context.Context, conversationID uuid.UUID) (persistence.Memory, error) {
	query := fmt.Sprintf(`SELECT id, uuid, conversation_id, query, response, created_at FROM %s WHERE conversation_id = $1`, r.tableName)
	row := r.db.QueryRowContext(ctx, query, conversationID)
	var memory persistence.Memory
	err := row.Scan(&memory.ID, &memory.UUID, &memory.ConversationID, &memory.Query, &memory.Response, &memory.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return persistence.Memory{}, fmt.Errorf("repo: memory not found for conversation id %s, %w", conversationID, err)
		}
		return persistence.Memory{}, fmt.Errorf("repo: error fetching memory for conversation id %s, %w", conversationID, err)
	}
	return memory, nil
}

func (r *MemoryRepo) FetchMany(ctx context.Context, memoryIds []int) ([]persistence.Memory, error) {
	if len(memoryIds) == 0 {
		return []persistence.Memory{}, nil
	}

	// Generate placeholders: $1, $2, $3, ...
	placeholders := make([]string, len(memoryIds))
	args := make([]interface{}, len(memoryIds))
	for i, id := range memoryIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, uuid, conversation_id, query, response, created_at FROM %s WHERE id IN (%s)`, r.tableName, strings.Join(placeholders, ", "))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repo: error fetching memories, %w", err)
	}
	defer rows.Close()
	var memories []persistence.Memory
	for rows.Next() {
		var memory persistence.Memory
		err := rows.Scan(&memory.ID, &memory.UUID, &memory.ConversationID, &memory.Query, &memory.Response, &memory.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("repo: error scanning memory, %w", err)
		}
		memories = append(memories, memory)
	}
	return memories, nil
}

func (r *MemoryRepo) FetchManyByConversationID(ctx context.Context, conversationID uuid.UUID, limit int) ([]persistence.Memory, error) {
	query := fmt.Sprintf(`SELECT id, uuid, conversation_id, query, response, created_at FROM %s WHERE conversation_id = $1`, r.tableName)
	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("repo: error fetching memories by conversation id %s, %w", conversationID, err)
	}
	var memories []persistence.Memory
	for rows.Next() {
		var memory persistence.Memory
		err = rows.Scan(&memory.ID, &memory.UUID, &memory.ConversationID, &memory.Query, &memory.Response, &memory.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("repo: error scanning memory, %w", err)
		}
		memories = append(memories, memory)
	}
	return memories, nil
}

func (r *MemoryRepo) InsertOne(ctx context.Context, conversationID uuid.UUID, memoryID uuid.UUID, query string, response string, createdAt time.Time) (int, error) {
	var insertedID int
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`INSERT INTO %s (conversation_id, uuid, query, response, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, r.tableName), conversationID, memoryID, query, response, createdAt).Scan(&insertedID)
	if err != nil {
		return 0, fmt.Errorf("repo: error inserting memory, %w", err)
	}
	return insertedID, nil
}

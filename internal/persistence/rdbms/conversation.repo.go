package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/haren7/minimal-memory/internal/persistence"

	"github.com/google/uuid"
)

type ConversationRepo struct {
	db *sql.DB
}

func NewConversationRepo(db *sql.DB) persistence.ConversationRepoInterface {
	return &ConversationRepo{db}
}

func (r *ConversationRepo) FetchOne(ctx context.Context, conversationID uuid.UUID) (persistence.Conversation, error) {
	query := "SELECT id, uuid, agent, user, created_at FROM conversations WHERE uuid = $1"
	row := r.db.QueryRowContext(ctx, query, conversationID)
	var conversation persistence.Conversation
	err := row.Scan(&conversation.ID, &conversation.UUID, &conversation.Agent, &conversation.User, &conversation.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return persistence.Conversation{}, fmt.Errorf("repo: conversation not found for id %s, %w", conversationID, err)
		}
		return persistence.Conversation{}, fmt.Errorf("repo: error fetching conversation for id %s, %w", conversationID, err)
	}
	return conversation, nil
}

func (r *ConversationRepo) InsertOne(ctx context.Context, agent string, user string, conversationID uuid.UUID, createdAt time.Time) (int, error) {
	var insertedID int
	err := r.db.QueryRowContext(ctx, "INSERT INTO conversations (uuid, agent, user, created_at) VALUES ($1, $2, $3, $4) RETURNING id", conversationID, agent, user, createdAt).Scan(&insertedID)
	if err != nil {
		return 0, fmt.Errorf("repo: error inserting conversation, %w", err)
	}
	return insertedID, nil
}

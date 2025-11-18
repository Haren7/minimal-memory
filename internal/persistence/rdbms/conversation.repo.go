package rdbms

import (
	"context"
	"database/sql"
	"memory/internal/persistence"
)

type ConversationRepo struct {
	db *sql.DB
}

func NewConversationRepo(db *sql.DB) persistence.RdbmsConversationRepoInterface {
	return &ConversationRepo{db}
}

func (r *ConversationRepo) FetchOne(ctx context.Context, agent string, user string) (persistence.RdbmsConversation, error) {
	return persistence.RdbmsConversation{}, nil
}

func (r *ConversationRepo) InsertOne(ctx context.Context, agent string, user string) (persistence.RdbmsConversation, error) {
	return persistence.RdbmsConversation{}, nil
}

package persistence

import (
	"time"

	"github.com/google/uuid"
)

type RdbmsConversation struct {
	ID        int       `db:"id"`
	UUID      uuid.UUID `db:"uuid"`
	Agent     string    `db:"agent"`
	User      string    `db:"user"`
	CreatedAt time.Time `db:"created_at"`
}

type RdbmsMemory struct {
	ID             int       `db:"id"`
	UUID           uuid.UUID `db:"uuid"`
	ConversationID uuid.UUID `db:"conversation_id"`
	Query          string    `db:"query"`
	Response       string    `db:"response"`
	CreatedAt      time.Time `db:"created_at"`
}

type VectorMemory struct {
	UUID      uuid.UUID
	Query     string
	Response  string
	CreatedAt time.Time
}

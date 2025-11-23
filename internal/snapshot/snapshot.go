package snapshot

import (
	"context"
)

type Manager interface {
	Store(ctx context.Context) error
	Load(ctx context.Context) error
}

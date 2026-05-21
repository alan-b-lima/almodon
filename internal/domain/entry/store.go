package entry

import (
	"context"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(ctx context.Context, order uuid.UUID) ([]Record, error)

	Get(ctx context.Context, order uuid.UUID, item uuid.UUID) (Record, error)

	Create(ctx context.Context, ent Entity) error

	Update(ctx context.Context, order uuid.UUID, item uuid.UUID, amount int64) error

	Delete(ctx context.Context, order uuid.UUID, item uuid.UUID) error
}

type (
	// FullRecord struct{}

	Record struct {
		Order  uuid.UUID
		Item   uuid.UUID
		Amount int64
	}
)

type (
	Entity struct {
		Order  uuid.UUID
		Item   uuid.UUID
		Amount int64
	}
)

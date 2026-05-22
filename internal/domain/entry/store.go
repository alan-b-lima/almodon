package entry

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(ctx context.Context, order uuid.UUID) ([]Record, error)

	Get(ctx context.Context, order uuid.UUID, item uuid.UUID) (Record, error)

	Create(ctx context.Context, order uuid.UUID, ent Entity) error

	Update(ctx context.Context, order uuid.UUID, item uuid.UUID, amount int64) error

	Delete(ctx context.Context, order uuid.UUID, item uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
	JoinTx(store.Store) (Store, error)
}

type (
	Record struct {
		Order  uuid.UUID
		Item   uuid.UUID
		Amount int64
	}
)

type (
	Entity struct {
		Item   uuid.UUID
		Amount int64
	}
)

package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	Get(ctx context.Context) (Record, error)

	Create(ctx context.Context, ent Entity) error

	Upgrade(ctx context.Context, order uuid.UUID) error

	RunTx(ctx context.Context, proc func(context.Context, Store) error) error
	JoinTx(store store.Store) (Store, error)
}

type (
	Record struct {
		UUID    uuid.UUID
		Bloom   uuid.UUID
		Version int
		Created time.Time
		Updated time.Time
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		Created time.Time
	}
)

package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(ctx context.Context) ([]Record, error)

	Get(ctx context.Context, uuid uuid.UUID) (Record, error)
	GetByName(ctx context.Context, name string) (Record, error)

	Create(ctx context.Context, ent Entity) error

	Rename(ctx context.Context, uuid uuid.UUID, name string) error
	Upgrade(ctx context.Context, stem uuid.UUID, order uuid.UUID) error

	Delete(ctx context.Context, uuid uuid.UUID) error

	RunTx(ctx context.Context, proc func(Store) error) error
	JoinTx(store store.Store) (Store, error)
}

type (
	Record struct {
		UUID    uuid.UUID
		Bloom   uuid.UUID
		Name    string
		Version int
		Created time.Time
		Updated time.Time
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		Name    string
		Created time.Time
	}
)

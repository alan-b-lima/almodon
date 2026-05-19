package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/order"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(ctx context.Context) ([]Record, error)

	Get(ctx context.Context, uuid uuid.UUID) (FullRecord, error)
	GetByTitle(ctx context.Context, title string) (FullRecord, error)

	Create(ctx context.Context, ent Entity) error

	Rename(ctx context.Context, uuid uuid.UUID, title string) error
	Upgrade(ctx context.Context, stem uuid.UUID, order uuid.UUID) error

	Delete(ctx context.Context, uuid uuid.UUID) error

	RunTx(ctx context.Context, proc func(Store) error) error
	JoinTx(store store.Store) (Store, error)
}

type (
	Record struct {
		UUID    uuid.UUID
		Bloom   uuid.UUID
		Title   string
		Version int
		Created time.Time
		Updated time.Time
	}

	FullRecord struct {
		UUID    uuid.UUID
		Bloom   uuid.UUID
		Title   string
		Version int
		Created time.Time
		Updated time.Time
		Orders  []order.Record
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		Bloom   uuid.UUID
		Title   string
		Created time.Time
	}
)

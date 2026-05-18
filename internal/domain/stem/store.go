package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/order"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(context.Context, uuid.UUID) ([]Record, error)

	Get(context.Context, uuid.UUID) (FullRecord, error)
	GetByTitle(context.Context, string) (FullRecord, error)

	Create(context.Context, Entity) error

	Rename(context.Context, uuid.UUID, string) error
	Upgrade(context.Context, uuid.UUID, uuid.UUID) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
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

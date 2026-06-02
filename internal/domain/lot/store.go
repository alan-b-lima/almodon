package lot

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context) ([]Record, error)
	ListByState(context.Context, State) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	Patch(context.Context, uuid.UUID, PatchEntity) error
	Sign(context.Context, uuid.UUID, SignEntity) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
	JoinTx(store.Store) (Store, error)
}

type (
	Record struct {
		UUID     uuid.UUID
		Order    uuid.UUID
		Supplier string
		Author   uuid.UUID
		Arrival  time.Time
		Note     string
		Created  time.Time
		Updated  time.Time
	}
)

type (
	Entity struct {
		UUID     uuid.UUID
		Supplier string
		Author   uuid.UUID
		Arrival  time.Time
		Note     string
		Created  time.Time
		Updated  time.Time
	}

	PatchEntity struct {
		Supplier opt.Opt[string]
		Arrival  opt.Opt[time.Time]
		Note     opt.Opt[string]
		Updated  time.Time
	}

	SignEntity struct {
		Order   uuid.UUID
		Updated time.Time
	}
)

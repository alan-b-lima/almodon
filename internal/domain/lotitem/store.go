package lotitem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context, uuid.UUID) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	Patch(context.Context, uuid.UUID, PatchEntity) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
	JoinTx(store.Store) (Store, error)
}

type (
	Record struct {
		UUID     uuid.UUID
		Lot      uuid.UUID
		Material uuid.UUID
		Name     string
		Unit     string
		Amount   int64
		UnitCost money.Money
		Expires  time.Time
		Created  time.Time
		Updated  time.Time
	}
)

type (
	Entity struct {
		UUID     uuid.UUID
		Lot      uuid.UUID
		Material uuid.UUID
		Amount   int64
		UnitCost money.Money
		Expires  time.Time
		Created  time.Time
		Updated  time.Time
	}

	PatchEntity struct {
		Material opt.Opt[uuid.UUID]
		Amount   opt.Opt[int64]
		UnitCost opt.Opt[money.Money]
		Expires  opt.Opt[time.Time]
		Updated  time.Time
	}
)

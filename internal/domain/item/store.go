package item

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context) ([]Record, error)
	ListByMaterial(context.Context, uuid.UUID) ([]Record, error)
	ListByECampus(context.Context, int) ([]Record, error)
	ListByCATMAT(context.Context, int) ([]Record, error)
	ListBySIADS(context.Context, int) ([]Record, error)
	ListByLot(context.Context, uuid.UUID) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	Update(context.Context, uuid.UUID, UpdateEntity) error
	Patch(context.Context, uuid.UUID, PatchEntity) error
	Effective(context.Context, EffectiveLot) error

	Delete(context.Context, uuid.UUID) error

	RunTx(context.Context, func(Store) error) error
	JoinTx(store.Store) (Store, error)
}

type (
	Record struct {
		UUID      uuid.UUID
		Material  uuid.UUID
		Name      string
		ECampus   int
		CATMAT    int
		SIADS     int
		Unit      string
		Lot       uuid.UUID
		Available int64
		UnitCost  money.Money
		Expires   time.Time
		Created   time.Time
		Updated   time.Time
	}
)

type (
	Entity struct {
		UUID      uuid.UUID
		Material  uuid.UUID
		Lot       uuid.UUID
		Available int64
		UnitCost  money.Money
		Expires   time.Time
		Created   time.Time
		Updated   time.Time
	}

	UpdateEntity struct {
		Amount  int64
		Updated time.Time
	}

	PatchEntity struct {
		UnitCost opt.Opt[money.Money]
		Expires  opt.Opt[time.Time]
		Updated  time.Time
	}

	EffectiveLot struct {
		Lot     uuid.UUID
		Updated time.Time
	}
)

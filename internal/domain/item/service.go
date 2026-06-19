package item

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(context.Context) ([]Result, error)
	ListByMaterial(context.Context, uuid.UUID) ([]Result, error)
	ListByECampus(context.Context, int) ([]Result, error)
	ListByCATMAT(context.Context, int) ([]Result, error)
	ListBySIADS(context.Context, int) ([]Result, error)
	ListByLot(context.Context, uuid.UUID) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)

	Create(context.Context, Create) (CreateResult, error)
	CreateMany(context.Context, []Create) error

	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Material uuid.UUID   `json:"material"`
		Lot      uuid.UUID   `json:"lot"`
		Original int64       `json:"original"`
		UnitCost money.Money `json:"unit_cost"`
		Expires  time.Time   `json:"expires"`
	}

	Patch struct {
		UnitCost opt.Opt[money.Money] `json:"unit_cost"`
		Expires  opt.Opt[time.Time]   `json:"expires"`
	}
)

type (
	Result struct {
		UUID        uuid.UUID   `json:"uuid"`
		Material    uuid.UUID   `json:"material"`
		Name        string      `json:"name"`
		ECampus     int         `json:"ecampus"`
		CATMAT      int         `json:"catmat"`
		SIADS       int         `json:"siads"`
		Unit        string      `json:"unit"`
		Lot         uuid.UUID   `json:"lot"`
		Available   int64       `json:"available"`
		UnitCost    money.Money `json:"unit_cost"`
		Expires     time.Time   `json:"expires"`
		ExpiresFlag Expiration  `json:"expires_flag"`
		Created     time.Time   `json:"created"`
		Updated     time.Time   `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

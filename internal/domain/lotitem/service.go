package lotitem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(ctx context.Context, lot uuid.UUID) ([]Result, error)

	Get(ctx context.Context, uuid uuid.UUID) (Result, error)

	Create(ctx context.Context, ent Create) (CreateResult, error)

	Patch(ctx context.Context, uuid uuid.UUID, ent Patch) error

	Delete(ctx context.Context, uuid uuid.UUID) error
}

type (
	Create struct {
		Lot      uuid.UUID   `json:"lot"`
		Material uuid.UUID   `json:"material"`
		Amount   int64       `json:"amount"`
		UnitCost money.Money `json:"unit_cost"`
		Expires  time.Time   `json:"expires"`
	}

	Patch struct {
		Material opt.Opt[uuid.UUID]   `json:"material"`
		Amount   opt.Opt[int64]       `json:"amount"`
		UnitCost opt.Opt[money.Money] `json:"unit_cost"`
		Expires  opt.Opt[time.Time]   `json:"expires"`
	}
)

type (
	Result struct {
		UUID     uuid.UUID   `json:"uuid"`
		Lot      uuid.UUID   `json:"lot"`
		Material uuid.UUID   `json:"material"`
		Name     string      `json:"name"`
		Unit     string      `json:"unit"`
		Amount   int64       `json:"amount"`
		UnitCost money.Money `json:"unit_cost"`
		Expires  time.Time   `json:"expires"`
		Created  time.Time   `json:"created"`
		Updated  time.Time   `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

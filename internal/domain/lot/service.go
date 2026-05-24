package lot

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(context.Context) ([]Result, error)
	ListByState(context.Context, State) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)

	Create(context.Context, Create) (CreateResult, error)

	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error

	Modify(context.Context, uuid.UUID, func(context.Context, Store) error) error
}

type (
	Create struct {
		Supplier string    `json:"suplier"`
		Arrival  time.Time `json:"arrival"`
		Note     string    `json:"note"`
	}

	Patch struct {
		Supplier opt.Opt[string]    `json:"suplier"`
		Arrival  opt.Opt[time.Time] `json:"arrival"`
		Note     opt.Opt[string]    `json:"note"`
	}
)

type (
	Result struct {
		UUID     uuid.UUID `json:"uuid"`
		Order    uuid.UUID `json:"order"`
		Supplier string    `json:"suplier"`
		Arrival  time.Time `json:"arrival"`
		Note     string    `json:"note"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

package material

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(context.Context) ([]Result, error)
	ListByECampus(context.Context, int) ([]Result, error)
	ListByCATMAT(context.Context, int) ([]Result, error)
	ListBySIADS(context.Context, int) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)

	Create(context.Context, Create) (CreateResult, error)

	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Name     string  `json:"name"`
		ECampus  int     `json:"ecampus"`
		CATMAT   int     `json:"catmat"`
		SIADS    int     `json:"siads"`
		Descript string  `json:"description"`
		Unit     string  `json:"unit"`
		Min      int64 `json:"min"`
	}

	Patch struct {
		Name     opt.Opt[string]  `json:"name"`
		ECampus  opt.Opt[int]     `json:"ecampus"`
		CATMAT   opt.Opt[int]     `json:"catmat"`
		SIADS    opt.Opt[int]     `json:"siads"`
		Descript opt.Opt[string]  `json:"description"`
		Unit     opt.Opt[string]  `json:"unit"`
		Min      opt.Opt[int64] `json:"min"`
	}
)

type (
	Result struct {
		UUID     uuid.UUID `json:"uuid"`
		Name     string    `json:"name"`
		ECampus  int       `json:"ecampus"`
		CATMAT   int       `json:"catmat"`
		SIADS    int       `json:"siads"`
		Descript string    `json:"description"`
		Unit     string    `json:"unit"`
		Min      int64   `json:"min"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

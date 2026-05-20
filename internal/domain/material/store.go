package material

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/opt"
)

type Store interface {
	List(context.Context) ([]Record, error)
	ListByECampus(context.Context, int) ([]Record, error)
	ListByCATMAT(context.Context, int) ([]Record, error)
	ListBySIADS(context.Context, int) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	Patch(context.Context, uuid.UUID, PatchEntity) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Record struct {
		UUID     uuid.UUID
		Name     string
		ECampus  int
		CATMAT   int
		SIADS    int
		Descript string
		Unit     string
		Min      float64
		Created  time.Time
		Updated  time.Time
	}

	Entity struct {
		UUID     uuid.UUID
		Name     string
		ECampus  int
		CATMAT   int
		SIADS    int
		Descript string
		Unit     string
		Min      float64
		Created  time.Time
		Updated  time.Time
	}

	PatchEntity struct {
		Name     opt.Opt[string]
		ECampus  opt.Opt[int]
		CATMAT   opt.Opt[int]
		SIADS    opt.Opt[int]
		Descript opt.Opt[string]
		Unit     opt.Opt[string]
		Min      opt.Opt[float64]
		Updated  time.Time
	}
)

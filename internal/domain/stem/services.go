package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(context.Context) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)
	GetByName(context.Context, string) (Result, error)

	Create(context.Context, Create) (CreateResult, error)

	Rename(context.Context, uuid.UUID, Rename) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Name string `json:"name"`
	}

	Rename struct {
		Name string `json:"name"`
	}
)

type (
	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		Bloom   uuid.UUID `json:"bloom"`
		Name    string    `json:"name"`
		Version int       `json:"version"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID    uuid.UUID `json:"uuid"`
		Version int       `json:"version"`
	}
)

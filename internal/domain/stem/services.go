package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/order"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	List(context.Context, uuid.UUID) ([]Result, error)

	Get(context.Context, uuid.UUID) (FullResult, error)
	GetByTitle(context.Context, string) (FullResult, error)

	Create(context.Context, Create) (CreateResult, error)

	Rename(context.Context, uuid.UUID, Rename) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		Title string `json:"title"`
	}

	Rename struct {
		Title string `json:"title"`
	}
)

type (
	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		Bloom   uuid.UUID `json:"bloom"`
		Title   string    `json:"title"`
		Version int       `json:"version"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}

	FullResult struct {
		UUID    uuid.UUID      `json:"uuid"`
		Bloom   uuid.UUID      `json:"bloom"`
		Title   string         `json:"title"`
		Version int            `json:"version"`
		Created time.Time      `json:"created"`
		Updated time.Time      `json:"updated"`
		Orders  []order.Result `json:"orders"`
	}

	CreateResult struct {
		UUID    uuid.UUID `json:"uuid"`
		Version int       `json:"version"`
	}
)

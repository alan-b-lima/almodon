package order

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	ListBlooms(ctx context.Context) ([]Result, error)
	ListByStem(ctx context.Context, stem uuid.UUID) ([]Result, error)

	Get(ctx context.Context, uuid uuid.UUID) (Result, error)
	GetBloom(ctx context.Context, stem uuid.UUID) (Result, error)

	Create(ctx context.Context, req Create) (CreateResult, error)
}

type (
	Create struct {
		Stem     uuid.UUID `json:"stem"`
		Password string    `json:"password"`
	}
)

type (
	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		Version int       `json:"version"`
		Stem    uuid.UUID `json:"stem"`
		Author  uuid.UUID `json:"author"`
		Name    string    `json:"name"`
		SIAPE   string    `json:"siape"`
		Created time.Time `json:"created"`
	}

	CreateResult struct {
		UUID    uuid.UUID `json:"uuid"`
		Version int       `json:"version"`
		Stem    uuid.UUID `json:"stem"`
	}
)

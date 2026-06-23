package order

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	ListBlooms(context.Context) ([]Record, error)
	ListByStem(context.Context, uuid.UUID) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)
	GetBloom(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		UUID    uuid.UUID
		Version int
		Stem    uuid.UUID
		Author  uuid.UUID
		Name    string
		SIAPE   string
		Created time.Time
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		Version int
		Stem    uuid.UUID
		Author  uuid.UUID
		Created time.Time
	}
)

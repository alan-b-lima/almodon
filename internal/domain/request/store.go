package request

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(context.Context) ([]Record, error)

	Get(context.Context, uuid.UUID) (Record, error)
	GetByNumber(context.Context, int) (Record, error)

	Create(context.Context, Entity) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Record struct {
		UUID    uuid.UUID
		Number  int
		Author  uuid.UUID
		Name    string
		SIAPE   string
		Title   string
		Memo    string
		Status  Status
		Created time.Time
		Updated time.Time
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		Number  int
		Author  uuid.UUID
		Title   string
		Memo    string
		Status  Status
		Created time.Time
		Updated time.Time
	}
)

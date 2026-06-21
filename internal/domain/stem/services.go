package stem

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Get(context.Context) (Result, error)

	Upgrade(context.Context, uuid.UUID) error
}

type (
	Result struct {
		UUID    uuid.UUID `json:"uuid"`
		Bloom   uuid.UUID `json:"bloom"`
		Version int       `json:"version"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
	}
)

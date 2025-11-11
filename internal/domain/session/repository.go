package session

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Getter
	Creater
	Updater
}

type (
	Getter interface {
		Get(uuid.UUID) (Entity, error)
	}

	Creater interface {
		Create(uuid.UUID, time.Duration) (Entity, error)
	}

	Updater interface {
		Update(uuid.UUID, time.Duration) (Entity, error)
	}
)

type (
	Entity struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}
)

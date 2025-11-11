package user

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Repository interface {
	Lister
	Getter
	GetterBySIAPE
	Creater
	Patcher
	Deleter
}

type (
	Lister interface {
		List(offset, limit int) (Entities, error)
	}

	Getter interface {
		Get(uuid uuid.UUID) (Entity, error)
	}

	GetterBySIAPE interface {
		GetBySIAPE(siape int) (Entity, error)
	}

	Creater interface {
		Create(siape int, name, email, password string, role auth.Role) (Entity, error)
	}

	Patcher interface {
		Patch(uuid uuid.UUID, name, email, password opt.Opt[string], role opt.Opt[auth.Role]) (Entity, error)
	}

	Deleter interface {
		Delete(uuid uuid.UUID) error
	}
)

type (
	Entities struct {
		Offset       int
		Length       int
		Records      []Entity
		TotalRecords int
	}

	Entity struct {
		UUID     uuid.UUID
		SIAPE    int
		Name     string
		Email    string
		Password [60]byte
		Role     auth.Role
	}

	AuthEntity struct {
		UUID    uuid.UUID
		User    uuid.UUID
		Expires time.Time
	}
)

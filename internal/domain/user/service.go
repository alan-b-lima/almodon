package user

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
)

type Service interface {
	List(context.Context) ([]Result, error)

	Get(context.Context, uuid.UUID) (Result, error)
	GetBySIAPE(context.Context, string) (Result, error)
	Me(context.Context) (Result, error)

	Auth(context.Context, uuid.UUID, string) error

	Create(context.Context, Create) (CreateResult, error)

	Patch(context.Context, uuid.UUID, Patch) error

	Delete(context.Context, uuid.UUID) error
}

type (
	Create struct {
		SIAPE    string    `json:"siape"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Role     auth.Role `json:"role"`
	}

	Patch struct {
		Name  opt.Opt[string] `json:"name"`
		Email opt.Opt[string] `json:"email"`
	}
)

type (
	Result struct {
		UUID     uuid.UUID `json:"uuid"`
		SIAPE    string    `json:"siape"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Role     auth.Role `json:"role"`
		Logged   bool      `json:"logged"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
	}

	CreateResult struct {
		UUID uuid.UUID `json:"uuid"`
	}
)

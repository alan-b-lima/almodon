package user

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type (
	ListRequest struct {
		Offset int `query:"offset"`
		Limit  int `query:"limit"`
	}

	GetRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	GetBySIAPERequest struct {
		SIAPE int `json:"-"`
	}

	CreateRequest struct {
		SIAPE    int    `json:"siape"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	PatchRequest struct {
		UUID     uuid.UUID       `json:"-"`
		Name     opt.Opt[string] `json:"name"`
		Email    opt.Opt[string] `json:"email"`
		Password opt.Opt[string] `json:"password"`
		Role     opt.Opt[string] `json:"role"`
	}

	DeleteRequest struct {
		UUID uuid.UUID `json:"-"`
	}

	AuthRequest struct {
		SIAPE    int    `json:"siape"`
		Password string `json:"password"`
	}

	ActorRequest struct {
		Session uuid.UUID `json:"-"`
	}
)

type (
	ListResponse struct {
		Offset       int        `json:"offset"`
		Length       int        `json:"length"`
		Records      []Response `json:"records"`
		TotalRecords int        `json:"total_records"`
	}

	Response struct {
		UUID  uuid.UUID `json:"uuid"`
		SIAPE int       `json:"siape"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Role  string    `json:"role"`
	}

	AuthResponse struct {
		UUID    uuid.UUID `json:"uuid"`
		User    uuid.UUID `json:"user"`
		Expires time.Time `json:"expires"`
	}
)

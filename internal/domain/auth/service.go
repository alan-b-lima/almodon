package auth

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Login(ctx context.Context, siape, password string) (Result, error)
	Logout(ctx context.Context, session session.Token) error

	Authenticator
}

type Authenticator interface {
	Actor(ctx context.Context, session session.Token) (Actor, error)
	Sign(ctx context.Context, user uuid.UUID, password string) error
}

type (
	Create struct {
		SIAPE    string `json:"siape"`
		Password string `json:"password"`
	}
)

type (
	Result struct {
		Token   session.Token `json:"-"`
		User    uuid.UUID     `json:"user"`
		Expires time.Time     `json:"expires"`
	}
)

package session

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	Get(context.Context, Token) (Result, error)

	Create(context.Context, Create) (Result, error)

	Update(context.Context, Token) error

	Delete(context.Context, Token) error

	ConfirmPassword(context.Context, Token) error
}

type (
	Create struct {
		User uuid.UUID `json:"user"`
	}
)

type (
	Result struct {
		Token            Token     `json:"-"`
		User             uuid.UUID `json:"user"`
		HardDeadline     time.Time `json:"hard_deadline"`
		IdleDeadline     time.Time `json:"idle_deadline"`
		PasswordVerified time.Time `json:"password_verified"`
	}
)

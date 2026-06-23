package session

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Store interface {
	List(context.Context) ([]Record, error)

	Get(context.Context, Token) (Record, error)
	GetByUser(context.Context, uuid.UUID) (Record, error)

	Create(context.Context, Entity) error

	UpdateActivity(context.Context, Token, time.Time) error
	UpdatePasswordVerified(context.Context, Token, time.Time) error

	Delete(context.Context, Token) error
	DeleteByUser(context.Context, uuid.UUID) error
	DeleteExpired(context.Context, time.Time) error

	RunTx(context.Context, func(Store) error) error
}

type (
	Record struct {
		Token            Token
		User             uuid.UUID
		HardDeadline     time.Time
		IdleDeadline     time.Time
		PasswordVerified time.Time
	}
)

type (
	Entity struct {
		Token            Token
		User             uuid.UUID
		HardDeadline     time.Time
		IdleDeadline     time.Time
		PasswordVerified time.Time
	}
)

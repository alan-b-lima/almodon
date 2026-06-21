package stemserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/stem"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	Service stem.Service
	Gate    auth.Authenticator
}

func NewGate(service stem.Service, gate auth.Authenticator) *Gate {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) Get(ctx context.Context) (stem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return stem.Result{}, err
	}

	return c.Service.Get(ctx)
}

func (c *Gate) Upgrade(ctx context.Context, order uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return err
	}

	return c.Service.Upgrade(ctx, order)
}

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
	stem.Service
	Gate auth.Authenticator
}

func NewGate(service stem.Service, gate auth.Authenticator) stem.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) List(ctx context.Context) ([]stem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return nil, err
	}

	return c.Service.List(ctx)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (stem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return stem.Result{}, err
	}

	return c.Service.Get(ctx, uuid)
}

func (c *Gate) GetByName(ctx context.Context, name string) (stem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return stem.Result{}, err
	}

	return c.Service.GetByName(ctx, name)
}

func (c *Gate) Create(ctx context.Context, req stem.Create) (stem.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return stem.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Rename(ctx context.Context, uuid uuid.UUID, req stem.Rename) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return err
	}

	return c.Service.Rename(ctx, uuid, req)
}

func (c *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockChange)
	if err != nil {
		return err
	}

	return c.Service.Delete(ctx, uuid)
}

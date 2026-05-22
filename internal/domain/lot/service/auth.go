package lotserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	lot.Service
	Gate auth.Authenticator
}

func NewGate(service lot.Service, gate auth.Authenticator) lot.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) List(ctx context.Context) ([]lot.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return nil, err
	}

	return c.Service.List(ctx)
}

func (c *Gate) ListByState(ctx context.Context, state lot.State) ([]lot.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByState(ctx, state)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (lot.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return lot.Result{}, err
	}

	return c.Service.Get(ctx, uuid)
}

func (c *Gate) Create(ctx context.Context, req lot.Create) (lot.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return lot.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Patch(ctx context.Context, uuid uuid.UUID, req lot.Patch) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return err
	}

	return c.Service.Patch(ctx, uuid, req)
}

func (c *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return err
	}

	return c.Service.Delete(ctx, uuid)
}

package lotitemserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/lotitem"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	lotitem.Service
	Gate auth.Authenticator
}

func NewGate(service lotitem.Service, gate auth.Authenticator) lotitem.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) List(ctx context.Context, lot uuid.UUID) ([]lotitem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return nil, err
	}

	return c.Service.List(ctx, lot)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (lotitem.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return lotitem.Result{}, err
	}

	return c.Service.Get(ctx, uuid)
}

func (c *Gate) Create(ctx context.Context, req lotitem.Create) (lotitem.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return lotitem.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Patch(ctx context.Context, uuid uuid.UUID, req lotitem.Patch) error {
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

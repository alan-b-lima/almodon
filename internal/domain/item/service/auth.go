package itemserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	item.Service
	Gate auth.Authenticator
}

func NewGate(service item.Service, gate auth.Authenticator) item.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) List(ctx context.Context) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.List(ctx)
}

func (c *Gate) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByMaterial(ctx, material)
}

func (c *Gate) ListByECampus(ctx context.Context, ecampus int) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByECampus(ctx, ecampus)
}

func (c *Gate) ListByCATMAT(ctx context.Context, catmat int) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByCATMAT(ctx, catmat)
}

func (c *Gate) ListBySIADS(ctx context.Context, siads int) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.ListBySIADS(ctx, siads)
}

func (c *Gate) ListByLot(ctx context.Context, lot uuid.UUID) ([]item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByLot(ctx, lot)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (item.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockRead)
	if err != nil {
		return item.Result{}, err
	}

	return c.Service.Get(ctx, uuid)
}

func (c *Gate) Create(ctx context.Context, req item.Create) (item.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return item.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Patch(ctx context.Context, uuid uuid.UUID, req item.Patch) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return err
	}

	return c.Service.Patch(ctx, uuid, req)
}

func (c *Gate) Effective(ctx context.Context, lot uuid.UUID) error {
	return auth.ErrBlocked
}

func (c *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.StockMgmt)
	if err != nil {
		return err
	}

	return c.Service.Delete(ctx, uuid)
}

package lotitemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/lotitem"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Items lotitem.Store
}

var _ lotitem.Service = (*Core)(nil)

func New(items lotitem.Store) *Core {
	return &Core{Items: items}
}

func (c *Core) List(ctx context.Context, lot uuid.UUID) ([]lotitem.Result, error) {
	recs, err := c.Items.List(ctx, lot)
	if err != nil {
		return nil, err
	}

	res := make([]lotitem.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, lotitem.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (lotitem.Result, error) {
	rec, err := c.Items.Get(ctx, uuid)
	if err != nil {
		return lotitem.Result{}, err
	}

	return lotitem.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req lotitem.Create) (lotitem.CreateResult, error) {
	ent := lotitem.Entity{
		Lot:      req.Lot,
		Material: req.Material,
	}

	if err := problem.Join(
		service.Set(&ent.Amount, req.Amount, lotitem.ProcessAmount),
		service.Set(&ent.UnitCost, req.UnitCost, lotitem.ProcessUnitCost),
		service.Set(&ent.Expires, req.Expires, lotitem.ProcessExpires),
	); err != nil {
		return lotitem.CreateResult{}, lotitem.ErrCreate.Cause(err).Make()
	}

	now := time.Now()
	ent.UUID = uuid.NewUUIDv7()
	ent.Created = now
	ent.Updated = now

	return lotitem.CreateResult{UUID: ent.UUID}, c.Items.Create(ctx, ent)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req lotitem.Patch) error {
	ent := lotitem.PatchEntity{Material: req.Material}
	if err := problem.Join(
		service.SetOpt(&ent.Amount, req.Amount, lotitem.ProcessAmount),
		service.SetOpt(&ent.UnitCost, req.UnitCost, lotitem.ProcessUnitCost),
		service.SetOpt(&ent.Expires, req.Expires, lotitem.ProcessExpires),
	); err != nil {
		return lotitem.ErrUpdate.Cause(err).Make()
	}

	ent.Updated = time.Now()

	return c.Items.Patch(ctx, uuid, ent)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Items.Delete(ctx, uuid)
}

package itemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Items item.Store
}

var _ item.Service = (*Core)(nil)

func New(items item.Store) *Core {
	return &Core{
		Items: items,
	}
}

func (c *Core) List(ctx context.Context) ([]item.Result, error) {
	recs, err := c.Items.List(ctx)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Result, error) {
	recs, err := c.Items.ListByMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) ListByECampus(ctx context.Context, ecampus int) ([]item.Result, error) {
	recs, err := c.Items.ListByECampus(ctx, ecampus)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) ListByCATMAT(ctx context.Context, catmat int) ([]item.Result, error) {
	recs, err := c.Items.ListByCATMAT(ctx, catmat)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) ListBySIADS(ctx context.Context, siads int) ([]item.Result, error) {
	recs, err := c.Items.ListBySIADS(ctx, siads)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) ListByLot(ctx context.Context, lot uuid.UUID) ([]item.Result, error) {
	recs, err := c.Items.ListByLot(ctx, lot)
	if err != nil {
		return nil, err
	}

	return translate_list(recs), nil
}

func (c *Core) Get(ctx context.Context, id uuid.UUID) (item.Result, error) {
	rec, err := c.Items.Get(ctx, id)
	if err != nil {
		return item.Result{}, err
	}

	return translate(&rec), nil
}

func (c *Core) Create(ctx context.Context, req item.Create) (item.CreateResult, error) {
	rec := item.Entity{
		Material: req.Material,
		Lot:      req.Lot,
	}

	err := problem.Join(
		service.Set(&rec.Available, req.Original, item.ProcessAmount),
		service.Set(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.Set(&rec.Expires, req.Expires, item.ProcessExpires),
	)
	if err != nil {
		return item.CreateResult{}, item.ErrCreate.Cause(err).Make()
	}

	now := time.Now()

	rec.UUID = uuid.NewUUIDv7()
	rec.Created = now
	rec.Updated = now

	return item.CreateResult{UUID: rec.UUID}, c.Items.Create(ctx, rec)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req item.Patch) error {
	var rec item.PatchEntity
	err := problem.Join(
		service.SetOpt(&rec.UnitCost, req.UnitCost, item.ProcessUnitCost),
		service.SetOpt(&rec.Expires, req.Expires, item.ProcessExpires),
	)
	if err != nil {
		return item.ErrUpdate.Cause(err).Make()
	}

	rec.Updated = time.Now()

	return c.Items.Patch(ctx, uuid, rec)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Items.Delete(ctx, uuid)
}

func translate(rec *item.Record) item.Result {
	return item.Result{
		UUID:        rec.UUID,
		Material:    rec.Material,
		Name:        rec.Name,
		ECampus:     rec.ECampus,
		CATMAT:      rec.CATMAT,
		SIADS:       rec.SIADS,
		Unit:        rec.Unit,
		Lot:         rec.Lot,
		Available:   rec.Available,
		UnitCost:    rec.UnitCost,
		Expires:     rec.Expires,
		ExpiresFlag: item.StatusExpires(rec.Expires),
		Created:     rec.Created,
		Updated:     rec.Updated,
	}
}

func translate_list(recs []item.Record) []item.Result {
	res := make([]item.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, translate(&rec))
	}

	return res
}

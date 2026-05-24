package lotserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Lots lot.Store
}

var _ lot.Service = (*Core)(nil)

func New(lots lot.Store) *Core {
	return &Core{Lots: lots}
}

func (c *Core) List(ctx context.Context) ([]lot.Result, error) {
	recs, err := c.Lots.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]lot.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, lot.Result(rec))
	}

	return res, nil
}

func (c *Core) ListByState(ctx context.Context, state lot.State) ([]lot.Result, error) {
	recs, err := c.Lots.ListByState(ctx, state)
	if err != nil {
		return nil, err
	}

	res := make([]lot.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, lot.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (lot.Result, error) {
	rec, err := c.Lots.Get(ctx, uuid)
	if err != nil {
		return lot.Result{}, err
	}

	return lot.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req lot.Create) (lot.CreateResult, error) {
	var ent lot.Entity
	if err := problem.Join(
		service.Set(&ent.Supplier, req.Supplier, lot.ProcessSupplier),
		service.Set(&ent.Arrival, req.Arrival, lot.ProcessArrival),
		service.Set(&ent.Note, req.Note, lot.ProcessNote),
	); err != nil {
		return lot.CreateResult{}, lot.ErrCreate.Cause(err).Make()
	}

	now := time.Now()
	ent.UUID = uuid.NewUUIDv7()
	ent.Created = now
	ent.Updated = now

	return lot.CreateResult{ent.UUID}, c.Lots.Create(ctx, ent)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req lot.Patch) error {
	var ent lot.PatchEntity
	if err := problem.Join(
		service.SetOpt(&ent.Supplier, req.Supplier, lot.ProcessSupplier),
		service.SetOpt(&ent.Arrival, req.Arrival, lot.ProcessArrival),
		service.SetOpt(&ent.Note, req.Note, lot.ProcessNote),
	); err != nil {
		return lot.ErrUpdate.Cause(err).Make()
	}

	ent.Updated = time.Now()

	return c.Modify(ctx, uuid, func(ctx context.Context, lots lot.Store) error {
		return lots.Patch(ctx, uuid, ent)
	})
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Modify(ctx, uuid, func(ctx context.Context, lots lot.Store) error {
		return lots.Delete(ctx, uuid)
	})
}

func (c *Core) Modify(ctx context.Context, uuid uuid.UUID, proc func(context.Context, lot.Store) error) error {
	return c.Lots.RunTx(ctx, func(lots lot.Store) error {
		rec, err := lots.Get(ctx, uuid)
		if err != nil {
			return err
		}

		if !rec.Order.IsZero() {
			return lot.ErrModifySigned
		}

		return proc(ctx, lots)
	})
}

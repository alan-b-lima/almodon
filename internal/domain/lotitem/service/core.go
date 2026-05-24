package lotitemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/domain/lotitem"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Items lotitem.Store
	Lots  lot.Service
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
	var ent lotitem.Entity
	if err := problem.Join(); err != nil {
		return lotitem.CreateResult{}, lotitem.ErrCreate.Cause(err).Make()
	}

	now := time.Now()
	ent.UUID = uuid.NewUUIDv7()
	ent.Created = now
	ent.Updated = now

	err := c.Lots.Modify(ctx, req.Lot, func(ctx context.Context, lots lot.Store) error {
		items, err := c.Items.JoinTx(lots)
		if err != nil {
			return err
		}

		return items.Create(ctx, ent)
	})
	return lotitem.CreateResult{ent.UUID}, err
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req lotitem.Patch) error {
	var ent lotitem.PatchEntity
	if err := problem.Join(); err != nil {
		return lotitem.ErrUpdate.Cause(err).Make()
	}

	ent.Updated = time.Now()

	rec, err := c.Items.Get(ctx, uuid)
	if err != nil {
		return err
	}

	return c.Lots.Modify(ctx, rec.Lot, func(ctx context.Context, lots lot.Store) error {
		items, err := c.Items.JoinTx(lots)
		if err != nil {
			return err
		}

		return items.Patch(ctx, uuid, ent)
	})
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	rec, err := c.Items.Get(ctx, uuid)
	if err != nil {
		return err
	}

	return c.Lots.Modify(ctx, rec.Lot, func(ctx context.Context, lots lot.Store) error {
		items, err := c.Items.JoinTx(lots)
		if err != nil {
			return err
		}

		return items.Delete(ctx, uuid)
	})
}

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
	Lots  lot.Store
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

	err := c.Items.RunTx(ctx, func(items lotitem.Store) error {
		lots, err := c.Lots.JoinTx(items)
		if err != nil {
			return err
		}

		if err := items.Create(ctx, ent); err != nil {
			return err
		}

		return lots.Mutable(ctx, req.Lot)
	})
	return lotitem.CreateResult{UUID: ent.UUID}, err
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

	return c.Items.RunTx(ctx, func(items lotitem.Store) error {
		lots, err := c.Lots.JoinTx(items)
		if err != nil {
			return err
		}

		if err := items.Patch(ctx, uuid, ent); err != nil {
			return err
		}

		return lots.Mutable(ctx, rec.Lot)
	})
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	rec, err := c.Items.Get(ctx, uuid)
	if err != nil {
		return err
	}

	return c.Items.RunTx(ctx, func(items lotitem.Store) error {
		lots, err := c.Lots.JoinTx(items)
		if err != nil {
			return err
		}

		if err := items.Delete(ctx, uuid); err != nil {
			return err
		}

		return lots.Mutable(ctx, rec.Lot)
	})
}

package stemserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/stem"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Stems stem.Store
}

var _ stem.Service = (*Core)(nil)

func New(stems stem.Store) *Core {
	return &Core{Stems: stems}
}

func (c *Core) List(ctx context.Context) ([]stem.Result, error) {
	recs, err := c.Stems.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]stem.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, stem.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (stem.Result, error) {
	rec, err := c.Stems.Get(ctx, uuid)
	if err != nil {
		return stem.Result{}, nil
	}

	return stem.Result(rec), nil
}

func (c *Core) GetByName(ctx context.Context, name string) (stem.Result, error) {
	rec, err := c.Stems.GetByName(ctx, name)
	if err != nil {
		return stem.Result{}, nil
	}

	return stem.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req stem.Create) (stem.CreateResult, error) {
	name, err := stem.ProcessName(req.Name)
	if err != nil {
		return stem.CreateResult{}, stem.ErrCreate.Cause(err).Make()
	}

	ent := stem.Entity{
		UUID:    uuid.NewUUIDv7(),
		Name:    name,
		Created: time.Now(),
	}

	return stem.CreateResult{UUID: ent.UUID}, c.Stems.Create(ctx, ent)
}

func (c *Core) Rename(ctx context.Context, uuid uuid.UUID, req stem.Rename) error {
	name, err := stem.ProcessName(req.Name)
	if err != nil {
		return stem.ErrCreate.Cause(err).Make()
	}

	return c.Stems.Rename(ctx, uuid, name)
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Stems.RunTx(ctx, func(c stem.Store) error {
		rec, err := c.Get(ctx, uuid)
		if err != nil {
			if err == stem.ErrNotFound {
				return nil
			}
			return err
		}

		if !rec.Bloom.IsZero() {
			return stem.ErrDeleteNonSprout
		}

		return c.Delete(ctx, uuid)
	})
}

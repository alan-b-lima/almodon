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

func (c *Core) Get(ctx context.Context) (stem.Result, error) {
	rec, err := c.Stems.Get(ctx)
	if err != nil {
		return stem.Result{}, nil
	}

	return stem.Result(rec), nil
}

func (c *Core) Create(ctx context.Context) (uuid.UUID, error) {
	ent := stem.Entity{
		UUID:    uuid.NewUUIDv7(),
		Created: time.Now(),
	}

	return ent.UUID, c.Stems.Create(ctx, ent)
}

func (c *Core) Upgrade(ctx context.Context, order uuid.UUID) error {
	return c.Stems.Upgrade(ctx, order)
}

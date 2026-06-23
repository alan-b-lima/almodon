package orderserve

import (
	"context"
	"errors"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/domain/order"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Core struct {
	Orders order.Store
	Lots   lot.Service
	Gate   auth.Authenticator
}

var _ order.Service = (*Core)(nil)

func New(orders order.Store, lots lot.Service, gate auth.Authenticator) *Core {
	return &Core{
		Orders: orders,
		Lots:   lots,
		Gate:   gate,
	}
}

func (c *Core) ListBlooms(ctx context.Context) ([]order.Result, error) {
	recs, err := c.Orders.ListBlooms(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]order.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, order.Result(rec))
	}

	return res, nil
}

func (c *Core) ListByStem(ctx context.Context, stem uuid.UUID) ([]order.Result, error) {
	recs, err := c.Orders.ListByStem(ctx, stem)
	if err != nil {
		return nil, err
	}

	res := make([]order.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, order.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (order.Result, error) {
	rec, err := c.Orders.Get(ctx, uuid)
	if err != nil {
		return order.Result{}, err
	}

	return order.Result(rec), nil
}

func (c *Core) GetBloom(ctx context.Context, stem uuid.UUID) (order.Result, error) {
	rec, err := c.Orders.GetBloom(ctx, stem)
	if err != nil {
		return order.Result{}, err
	}

	return order.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req order.Create) (order.CreateResult, error) {
	author, err := c.author(ctx, req.Password)
	if err != nil {
		return order.CreateResult{}, err
	}

	var res order.CreateResult
	err = c.Orders.RunTx(ctx, func(orders order.Store) error {
		rec, err := orders.GetBloom(ctx, req.Stem)
		if err != nil {
			if err != order.ErrNotFound {
				return err
			}

			rec = order.Record{
				Version: 0,
				Stem:    req.Stem,
			}
		}

		ent := order.Entity{
			UUID:    uuid.NewUUIDv7(),
			Version: order.NextVersion(rec.Version),
			Stem:    rec.Stem,
			Author:  author,
			Created: time.Now(),
		}

		if err := orders.Create(ctx, ent); err != nil {
			return err
		}

		res = order.CreateResult{
			UUID:    ent.UUID,
			Version: ent.Version,
			Stem:    ent.Stem,
		}
		return nil
	})
	return res, err
}

func (c *Core) author(ctx context.Context, password string) (uuid.UUID, error) {
	actor, err := service.ActorFromContext(ctx, c.Gate)
	if err != nil {
		return uuid.UUID{}, err
	}

	if err := c.Gate.Sign(ctx, actor.User, password); err != nil {
		if errors.Is(err, auth.ErrUnauthorized.Make()) {
			return uuid.UUID{}, order.ErrSignatureFailed
		}

		return uuid.UUID{}, err
	}

	return actor.User, nil
}

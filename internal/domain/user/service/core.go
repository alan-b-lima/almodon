package userserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Users user.Store
}

var _ user.Service = (*Core)(nil)

func New(users user.Store) *Core {
	return &Core{
		Users: users,
	}
}

func (c *Core) List(ctx context.Context) ([]user.Result, error) {
	recs, err := c.Users.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]user.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, translate(&rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, uuid uuid.UUID) (user.Result, error) {
	rec, err := c.Users.Get(ctx, uuid)
	if err != nil {
		return user.Result{}, err
	}

	return translate(&rec), nil
}

func (c *Core) GetBySIAPE(ctx context.Context, siape string) (user.Result, error) {
	rec, err := c.Users.GetBySIAPE(ctx, siape)
	if err != nil {
		return user.Result{}, err
	}

	return translate(&rec), nil
}

func (c *Core) Me(ctx context.Context) (user.Result, error) {
	return user.Result{}, support.ErrTODO
}

func (c *Core) Auth(ctx context.Context, uuid uuid.UUID, password string) error {
	rec, err := c.Users.Get(ctx, uuid)
	if err != nil {
		return err
	}

	return user.ComparePassword(rec.Password, password)
}

func (c *Core) Create(ctx context.Context, req user.Create) (user.CreateResult, error) {
	var rec user.Entity
	err := problem.Join(
		service.Set(&rec.SIAPE, req.SIAPE, user.ProcessSIAPE),
		service.Set(&rec.Name, req.Name, user.ProcessName),
		service.Set(&rec.Email, req.Email, user.ProcessEmail),
		service.Set(&rec.Password, req.Password, user.ProcessPassword),
		service.Set(&rec.Role, req.Role, user.ProcessRole),
	)
	if err != nil {
		return user.CreateResult{}, user.ErrCreate.Cause(err).Make()
	}

	now := time.Now()

	rec.UUID = uuid.NewUUIDv7()
	rec.Created = now
	rec.Updated = now

	return user.CreateResult{UUID: rec.UUID}, c.Users.Create(ctx, rec)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req user.Patch) error {
	var rec user.PatchEntity
	err := problem.Join(
		service.SetOpt(&rec.Name, req.Name, user.ProcessName),
		service.SetOpt(&rec.Email, req.Email, user.ProcessEmail),
	)
	if err != nil {
		return user.ErrUpdate.Cause(err).Make()
	}

	rec.Updated = time.Now()

	return c.Users.RunTx(ctx, func(c user.Store) error {
		if err := c.Patch(ctx, uuid, rec); err != nil {
			return err
		}

		count, err := c.CountChiefs(ctx)
		if err != nil {
			return err
		}

		if count <= 0 {
			return user.ErrNotEnoughChiefs
		}
		return nil
	})
}

func (c *Core) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.Users.RunTx(ctx, func(c user.Store) error {
		if err := c.Delete(ctx, uuid); err != nil {
			return err
		}

		count, err := c.CountChiefs(ctx)
		if err != nil {
			return err
		}

		if count <= 0 {
			return user.ErrNotEnoughChiefs
		}
		return nil
	})
}

func translate(rec *user.Record) user.Result {
	return user.Result{
		UUID:    rec.UUID,
		SIAPE:   rec.SIAPE,
		Name:    rec.Name,
		Email:   rec.Email,
		Role:    rec.Role,
		Logged:  rec.Logged,
		Created: rec.Created,
		Updated: rec.Updated,
	}
}

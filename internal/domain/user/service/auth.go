package userserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	user.Service
	Gate auth.Authenticator
}

func NewGate(service user.Service, gate auth.Authenticator) user.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

func (c *Gate) List(ctx context.Context) ([]user.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		return []user.Result{}, err
	}

	return c.Service.List(ctx)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (user.Result, error) {
	ctx, actor, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		if actor.User == uuid {
			goto Do
		}

		return user.Result{}, err
	}

Do:
	return c.Service.Get(ctx, uuid)
}

func (c *Gate) GetBySIAPE(ctx context.Context, siape string) (user.Result, error) {
	res, err := c.Service.GetBySIAPE(ctx, siape)
	if err != nil {
		return user.Result{}, err
	}

	ctx, actor, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		if actor.User == res.UUID {
			goto Do
		}

		return user.Result{}, err
	}

Do:
	return res, nil
}

func (c *Gate) Me(ctx context.Context) (user.Result, error) {
	ctx, actor, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserSelf)
	if err != nil {
		return user.Result{}, err
	}

	return c.Service.Get(ctx, actor.User)
}

func (c *Gate) Create(ctx context.Context, req user.Create) (user.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		return user.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Patch(ctx context.Context, uuid uuid.UUID, req user.Patch) error {
	ctx, actor, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		if actor.User == uuid {
			goto Do
		}

		return err
	}

Do:
	return c.Service.Patch(ctx, uuid, req)
}

func (c *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, actor, err := service.AuthorizeFromContext(ctx, c.Gate, perms.UserMgmt)
	if err != nil {
		if actor.User == uuid {
			goto Do
		}

		return err
	}

Do:
	return c.Service.Delete(ctx, uuid)
}

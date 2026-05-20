package materialserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	material.Service
	Gate auth.Authenticator
}

func NewGate(service material.Service, gate auth.Authenticator) material.Service {
	return &Gate{
		Service: service,
		Gate:    gate,
	}
}

var (
	perm_admin = auth.Allow(auth.Admin)
	perm_user  = auth.Allow(auth.User)
)

func (c *Gate) List(ctx context.Context) ([]material.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_user)
	if err != nil {
		return nil, err
	}

	return c.Service.List(ctx)
}

func (c *Gate) ListByECampus(ctx context.Context, ecampus int) ([]material.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_user)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByECampus(ctx, ecampus)
}

func (c *Gate) ListByCATMAT(ctx context.Context, catmat int) ([]material.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_user)
	if err != nil {
		return nil, err
	}

	return c.Service.ListByCATMAT(ctx, catmat)
}

func (c *Gate) ListBySIADS(ctx context.Context, siads int) ([]material.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_user)
	if err != nil {
		return nil, err
	}

	return c.Service.ListBySIADS(ctx, siads)
}

func (c *Gate) Get(ctx context.Context, uuid uuid.UUID) (material.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_user)
	if err != nil {
		return material.Result{}, err
	}

	return c.Service.Get(ctx, uuid)
}

func (c *Gate) Create(ctx context.Context, req material.Create) (material.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_admin)
	if err != nil {
		return material.CreateResult{}, err
	}

	return c.Service.Create(ctx, req)
}

func (c *Gate) Patch(ctx context.Context, uuid uuid.UUID, req material.Patch) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_admin)
	if err != nil {
		return err
	}

	return c.Service.Patch(ctx, uuid, req)
}

func (c *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, c.Gate, perm_admin)
	if err != nil {
		return err
	}

	return c.Service.Delete(ctx, uuid)
}

package promotionserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/auth/perms"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Gate struct {
	promotion.Service
	Gate auth.Authenticator
}

func NewGate(promotions promotion.Service, gate auth.Authenticator) promotion.Service {
	return &Gate{
		Service: promotions,
		Gate:    gate,
	}
}

func (g *Gate) Get(ctx context.Context, uuid uuid.UUID) (promotion.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, g.Gate, perms.UserMgmt)
	if err != nil {
		return promotion.Result{}, err
	}

	return g.Service.Get(ctx, uuid)
}

func (g *Gate) GetByUser(ctx context.Context, uuid uuid.UUID) (promotion.Result, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, g.Gate, perms.UserMgmt)
	if err != nil {
		return promotion.Result{}, err
	}

	return g.Service.GetByUser(ctx, uuid)
}

func (g *Gate) Create(ctx context.Context, req promotion.Create) (promotion.CreateResult, error) {
	ctx, _, err := service.AuthorizeFromContext(ctx, g.Gate, perms.UserMgmt)
	if err != nil {
		return promotion.CreateResult{}, err
	}

	return g.Service.Create(ctx, req)
}

func (g *Gate) Update(ctx context.Context, uuid uuid.UUID, req promotion.Update) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, g.Gate, perms.UserMgmt)
	if err != nil {
		return err
	}

	return g.Service.Update(ctx, uuid, req)
}

func (g *Gate) Delete(ctx context.Context, uuid uuid.UUID) error {
	ctx, _, err := service.AuthorizeFromContext(ctx, g.Gate, perms.UserMgmt)
	if err != nil {
		return err
	}

	return g.Service.Delete(ctx, uuid)
}

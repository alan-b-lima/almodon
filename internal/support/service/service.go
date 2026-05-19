package service

import (
	"context"
	"errors"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/session"

	"github.com/alan-b-lima/almodon/pkg/rbac"
)

func AuthorizeFromContext(ctx context.Context, gate auth.Authenticator, perms rbac.Permission[auth.Role]) (context.Context, auth.Actor, error) {
	actor, err := ActorFromContext(ctx, gate)
	if err != nil {
		ctx = context.WithValue(ctx, "actor", result{Err: err})
		return ctx, auth.Actor{}, err
	}

	if role := actor.Role; !perms.Allows(role) {
		return ctx, auth.Actor{}, auth.ErrUnauthorized.Details(map[string]any{"allowed": perms}).Make(role, perms)
	}

	ctx = context.WithValue(ctx, "actor", result{Actor: actor})
	return ctx, actor, nil
}

type result struct {
	auth.Actor
	Err error
}

func ActorFromContext(ctx context.Context, gate auth.Authenticator) (auth.Actor, error) {
	if actor, ok := ctx.Value("actor").(result); ok {
		return actor.Actor, actor.Err
	}

	session, ok := ctx.Value("session").(session.Token)
	if !ok {
		return auth.NewUnlogged(), nil
	}

	actor, err := gate.Actor(ctx, session)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthenticated.Make()) {
			return auth.Actor{}, err
		}

		return auth.NewUnlogged(), nil
	}

	return actor, nil
}

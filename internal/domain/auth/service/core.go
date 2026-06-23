package authserve

import (
	"context"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/domain/user"
)

type Core struct {
	Users    user.Service
	Sessions session.Service
}

var _ auth.Service = &Core{}

func New(users user.Service, sessions session.Service) *Core {
	return &Core{
		Users:    users,
		Sessions: sessions,
	}
}

func (c *Core) Login(ctx context.Context, siape string, password string) (auth.Result, error) {
	res, err := c.Users.GetBySIAPE(ctx, siape)
	if err != nil {
		return auth.Result{}, err
	}

	if err := user.ComparePassword(res.Password, password); err != nil {
		return auth.Result{}, err
	}

	sres, err := c.Sessions.Create(ctx, session.Create{User: res.UUID})
	if err != nil {
		return auth.Result{}, err
	}

	return auth.Result{
		Token:   sres.Token,
		User:    sres.User,
		Expires: sres.HardDeadline,
	}, nil
}

func (c *Core) Logout(ctx context.Context, session session.Token) error {
	return c.Sessions.Delete(ctx, session)
}

func (c *Core) Actor(ctx context.Context, session session.Token) (auth.Actor, error) {
	res, err := c.Sessions.Get(ctx, session)
	if err != nil {
		return auth.NewUnlogged(), auth.ErrUnauthenticated.Cause(err).Make()
	}

	ures, err := c.Users.Get(ctx, res.User)
	if err != nil {
		return auth.NewUnlogged(), auth.ErrUnauthenticated.Cause(err).Make()
	}

	return auth.NewLogged(ures.UUID, ures.Role), nil
}

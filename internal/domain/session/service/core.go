package sessionserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"

	"github.com/alan-b-lima/pkg/scheduler"
)

type Core struct {
	Sessions  session.Store
	Scheduler *scheduler.Scheduler
}

var _ session.Service = (*Core)(nil)

func New(sessions session.Store, scheduler *scheduler.Scheduler) *Core {
	return &Core{
		Sessions:  sessions,
		Scheduler: scheduler,
	}
}

func (c *Core) Get(ctx context.Context, token session.Token) (session.Result, error) {
	rec, err := c.Sessions.Get(ctx, token)
	if err != nil {
		return session.Result{}, err
	}

	if session.Expired(rec.HardDeadline, rec.IdleDeadline) {
		return session.Result{}, session.ErrNotFound
	}

	return session.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req session.Create) (session.Result, error) {
	rec := session.Entity{
		Token:        session.NewToken(),
		User:         req.User,
		HardDeadline: time.Now().Add(session.HardTimeout),
		IdleDeadline: time.Now().Add(session.IdleTimeout),
		PasswordVerified: time.Now(), 
	}

	err := c.Sessions.RunTx(ctx, func(store session.Store) error {
		err := store.DeleteByUser(ctx, rec.User)
		if err != nil {
			return err
		}
		return store.Create(ctx, rec)
	})

	if err != nil {
		return session.Result{}, err
	}

	c.scheduleCleanup(rec.HardDeadline)

	return session.Result(rec), nil
}

func (c *Core) ConfirmPassword(ctx context.Context, token session.Token) error {
	if _, err := c.Get(ctx, token); err != nil {
		return err
	}

	return c.Sessions.UpdatePasswordVerified(ctx, token, time.Now())
}

func (c *Core) Update(ctx context.Context, token session.Token) error {
	rec, err := c.Sessions.Get(ctx, token)
	if err != nil {
		return err
	}

	if session.Expired(rec.HardDeadline, rec.IdleDeadline) {
		return session.ErrNotFound
	}

	updateRec := time.Now().Add(session.IdleTimeout)

	if err := c.Sessions.UpdateActivity(ctx, token, updateRec); err != nil {
		return err
	}

	return nil
}

func (c *Core) Delete(ctx context.Context, token session.Token) error {
	return c.Sessions.Delete(ctx, token)
}

func (c *Core) Publish(ctx context.Context) error {
	err := c.Sessions.DeleteExpired(ctx, time.Now())
	if err != nil {
		return err
	}

	recs, err := c.Sessions.List(ctx)
	if err != nil {
		return err
	}

	for _, rec := range recs {
		c.scheduleCleanup(rec.HardDeadline)
	}

	return nil
}

func (c *Core) scheduleCleanup(expiresAt time.Time) {
	c.Scheduler.Post(func() {
		_ = c.Sessions.DeleteExpired(context.Background(), expiresAt)
	}, expiresAt)
}

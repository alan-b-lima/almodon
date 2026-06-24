package sessionserve

import (
	"cmp"
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"

	"github.com/alan-b-lima/pkg/scheduler"
)

type Core struct {
	sessions  session.Store
	scheduler *scheduler.Scheduler

	hard_deadline time.Duration
	idle_deadline time.Duration
	pure_get      bool
}

type Opts struct {
	HardDeadline time.Duration
	IdleDeadline time.Duration
	PureGet      bool
}

var _ session.Service = (*Core)(nil)

func New(sessions session.Store, scheduler *scheduler.Scheduler, opts *Opts) *Core {
	return &Core{
		sessions:      sessions,
		scheduler:     scheduler,
		hard_deadline: cmp.Or(opts.HardDeadline, session.DefaultHardTimeout),
		idle_deadline: cmp.Or(opts.IdleDeadline, session.DefaultIdleTimeout),
		pure_get:      opts.PureGet,
	}
}

func (c *Core) Get(ctx context.Context, token session.Token) (session.Result, error) {
	if c.pure_get {
		return c.GetPure(ctx, token)
	}

	var res session.Result
	err := c.sessions.RunTx(ctx, func(store session.Store) (err error) {
		c := c.extend(store)

		if res, err = c.GetPure(ctx, token); err != nil {
			return err
		}

		return c.Update(ctx, token)
	})
	if err != nil {
		return session.Result{}, err
	}

	return res, nil
}

func (c *Core) GetPure(ctx context.Context, token session.Token) (session.Result, error) {
	rec, err := c.sessions.Get(ctx, token)
	if err != nil {
		return session.Result{}, err
	}

	if session.IsExpired(rec.HardDeadline, rec.IdleDeadline) {
		return session.Result{}, session.ErrNotFound
	}

	return session.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req session.Create) (session.Result, error) {
	now := time.Now()

	rec := session.Entity{
		Token:        session.NewToken(),
		User:         req.User,
		HardDeadline: now.Add(session.DefaultHardTimeout),
		IdleDeadline: now.Add(session.DefaultIdleTimeout),
	}

	err := c.sessions.RunTx(ctx, func(store session.Store) error {
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

func (c *Core) Update(ctx context.Context, token session.Token) error {
	deadline := time.Now().Add(session.DefaultIdleTimeout)
	return c.sessions.Update(ctx, token, deadline)
}

func (c *Core) Delete(ctx context.Context, token session.Token) error {
	return c.sessions.Delete(ctx, token)
}

func (c *Core) extend(store session.Store) *Core {
	cc := *c
	cc.sessions = store
	return &cc
}

func (c *Core) Publish(ctx context.Context) error {
	err := c.sessions.DeleteExpired(ctx, time.Now())
	if err != nil {
		return err
	}

	recs, err := c.sessions.List(ctx)
	if err != nil {
		return err
	}

	for _, rec := range recs {
		c.scheduleCleanup(rec.HardDeadline)
	}

	return nil
}

func (c *Core) scheduleCleanup(expires time.Time) {
	c.scheduler.Post(func() {
		_ = c.sessions.DeleteExpired(context.TODO(), expires)
	}, expires)
}

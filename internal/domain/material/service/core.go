package materialserve

import (
	"context"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/material"
	"github.com/alan-b-lima/almodon/internal/support/service"

	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

type Core struct {
	Materials material.Store
}

var _ material.Service = (*Core)(nil)

func New(materials material.Store) *Core {
	return &Core{
		Materials: materials,
	}
}

func (c *Core) List(ctx context.Context) ([]material.Result, error) {
	recs, err := c.Materials.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]material.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, material.Result(rec))
	}

	return res, nil
}

func (c *Core) ListByECampus(ctx context.Context, ecampus int) ([]material.Result, error) {
	recs, err := c.Materials.ListByECampus(ctx, ecampus)
	if err != nil {
		return nil, err
	}

	res := make([]material.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, material.Result(rec))
	}

	return res, nil
}

func (c *Core) ListByCATMAT(ctx context.Context, catmat int) ([]material.Result, error) {
	recs, err := c.Materials.ListByCATMAT(ctx, catmat)
	if err != nil {
		return nil, err
	}

	res := make([]material.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, material.Result(rec))
	}

	return res, nil
}

func (c *Core) ListBySIADS(ctx context.Context, siads int) ([]material.Result, error) {
	recs, err := c.Materials.ListBySIADS(ctx, siads)
	if err != nil {
		return nil, err
	}

	res := make([]material.Result, 0, len(recs))
	for _, rec := range recs {
		res = append(res, material.Result(rec))
	}

	return res, nil
}

func (c *Core) Get(ctx context.Context, id uuid.UUID) (material.Result, error) {
	rec, err := c.Materials.Get(ctx, id)
	if err != nil {
		return material.Result{}, err
	}

	return material.Result(rec), nil
}

func (c *Core) Create(ctx context.Context, req material.Create) (material.CreateResult, error) {
	var rec material.Entity
	err := problem.Join(
		service.Set(&rec.Name, req.Name, material.ProcessName),
		service.Set(&rec.ECampus, req.ECampus, material.ProcessECampus),
		service.Set(&rec.CATMAT, req.CATMAT, material.ProcessCATMAT),
		service.Set(&rec.SIADS, req.SIADS, material.ProcessSIADS),
		service.Set(&rec.Descript, req.Descript, material.ProcessDescription),
		service.Set(&rec.Unit, req.Unit, material.ProcessUnit),
		service.Set(&rec.Min, req.Min, material.ProcessMin),
	)
	if err != nil {
		return material.CreateResult{}, material.ErrCreate.Cause(err).Make()
	}

	now := time.Now()
	rec.UUID = uuid.NewUUIDv7()
	rec.Created = now
	rec.Updated = now

	return material.CreateResult{UUID: rec.UUID}, c.Materials.Create(ctx, rec)
}

func (c *Core) Patch(ctx context.Context, uuid uuid.UUID, req material.Patch) error {
	var rec material.PatchEntity
	err := problem.Join(
		service.SetOpt(&rec.Name, req.Name, material.ProcessName),
		service.SetOpt(&rec.ECampus, req.ECampus, material.ProcessECampus),
		service.SetOpt(&rec.CATMAT, req.CATMAT, material.ProcessCATMAT),
		service.SetOpt(&rec.SIADS, req.SIADS, material.ProcessSIADS),
		service.SetOpt(&rec.Descript, req.Descript, material.ProcessDescription),
		service.SetOpt(&rec.Unit, req.Unit, material.ProcessUnit),
		service.SetOpt(&rec.Min, req.Min, material.ProcessMin),
	)
	if err != nil {
		return material.ErrUpdate.Cause(err).Make()
	}

	rec.Updated = time.Now()

	return c.Materials.Patch(ctx, uuid, rec)
}

func (c *Core) Delete(ctx context.Context, id uuid.UUID) error {
	return c.Materials.Delete(ctx, id)
}

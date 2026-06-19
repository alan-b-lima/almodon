package items

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/resource"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux

	Items item.Service
}

func New(items item.Service) *Resource {
	rc := Resource{
		Items: items,
	}

	routes := map[string]http.HandlerFunc{
		"GET /items/{$}":                 rc.List,
		"GET /items/material/{material}": rc.ListByMaterial,
		"GET /items/ecampus/{ecampus}":   rc.ListByECampus,
		"GET /items/catmat/{catmat}":     rc.ListByCATMAT,
		"GET /items/siads/{siads}":       rc.ListBySIADS,
		"GET /items/{uuid}":              rc.Get,
		"PATCH /items/{uuid}":            rc.Patch,
		"DELETE /items/{uuid}":           rc.Delete,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Items.List, w, r)
}

func (rc *Resource) ListByMaterial(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]item.Result, error) {
		material, err := uuid.FromString(r.PathValue("material"))
		if err != nil {
			return nil, resource.ErrBadUUID
		}

		return rc.Items.ListByMaterial(ctx, material)
	}, w, r)
}

func (rc *Resource) ListByECampus(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]item.Result, error) {
		ecampus, err := strconv.Atoi(r.PathValue("ecampus"))
		if err != nil {
			return nil, resource.ErrBadInteger
		}

		return rc.Items.ListByECampus(ctx, ecampus)
	}, w, r)
}

func (rc *Resource) ListByCATMAT(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]item.Result, error) {
		catmat, err := strconv.Atoi(r.PathValue("catmat"))
		if err != nil {
			return nil, resource.ErrBadInteger
		}

		return rc.Items.ListByCATMAT(ctx, catmat)
	}, w, r)
}

func (rc *Resource) ListBySIADS(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]item.Result, error) {
		siads, err := strconv.Atoi(r.PathValue("siads"))
		if err != nil {
			return nil, resource.ErrBadInteger
		}

		return rc.Items.ListBySIADS(ctx, siads)
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (item.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return item.Result{}, resource.ErrBadUUID
		}

		return rc.Items.Get(ctx, uuid)
	}, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req item.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Items.Patch(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Items.Delete(ctx, uuid)
	}, w, r)
}

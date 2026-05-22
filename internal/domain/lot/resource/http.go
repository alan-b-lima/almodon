package lots

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/lot"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Lots lot.Service
}

func New(lots lot.Service) *Resource {
	rc := Resource{
		Lots: lots,
	}

	routes := map[string]http.HandlerFunc{
		"GET /lots/{state}":   rc.List,
		"GET /lots/{uuid}":    rc.Get,
		"POST /lots/{$}":      rc.Create,
		"PATCH /lots/{uuid}":  rc.Patch,
		"DELETE /lots/{uuid}": rc.Delete,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]lot.Result, error) {
		switch r.PathValue("state") {
		case "":
			return rc.Lots.List(ctx)
		case "open":
			return rc.Lots.ListByState(ctx, lot.Open)
		case "closed":
			return rc.Lots.ListByState(ctx, lot.Closed)
		}

		return nil, resource.ErrNotFound.Make(r.URL.Path)
	}, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (lot.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return lot.Result{}, resource.ErrBadUUID
		}

		return rc.Lots.Get(ctx, uuid)
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), rc.Lots.Create, w, r)
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req lot.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Lots.Patch(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Lots.Delete(ctx, uuid)
	}, w, r)
}

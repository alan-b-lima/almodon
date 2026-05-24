package lots

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/lot"
	item "github.com/alan-b-lima/almodon/internal/domain/lotitem"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Lots  lot.Service
	Items item.Service
}

func New(lots lot.Service, items item.Service) *Resource {
	rc := Resource{
		Lots:  lots,
		Items: items,
	}

	routes := map[string]http.HandlerFunc{
		"GET /lots/":                 rc.List,
		"GET /lots/open":             rc.ListOpen,
		"GET /lots/closed":           rc.ListClosed,
		"GET /lots/{uuid}":           rc.Get,
		"POST /lots/{$}":             rc.Create,
		"PATCH /lots/{uuid}":         rc.Patch,
		"DELETE /lots/{uuid}":        rc.Delete,
		"GET /lots/{uuid}/items/{$}": rc.ListItems,
		"GET /lots/items/{uuid}":     rc.GetItems,
		"POST /lots/items/{uuid}":    rc.CreateItems,
		"PATCH /lots/items/{uuid}":   rc.PatchItems,
		"DELETE /lots/items/{uuid}":  rc.DeleteItems,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Lots.List, w, r)
}

func (rc *Resource) ListOpen(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]lot.Result, error) {
		return rc.Lots.ListByState(ctx, lot.Open)
	}, w, r)
}

func (rc *Resource) ListClosed(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]lot.Result, error) {
		return rc.Lots.ListByState(ctx, lot.Closed)
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

func (rc *Resource) ListItems(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) ([]item.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return nil, resource.ErrBadUUID
		}

		return rc.Items.List(ctx, uuid)
	}, w, r)
}

func (rc *Resource) GetItems(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (item.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return item.Result{}, resource.ErrBadUUID
		}

		return rc.Items.Get(ctx, uuid)
	}, w, r)
}

func (rc *Resource) CreateItems(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), rc.Items.Create, w, r)
}

func (rc *Resource) PatchItems(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, ent item.Patch) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Items.Patch(ctx, uuid, ent)
	}, w, r)
}

func (rc *Resource) DeleteItems(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Items.Delete(ctx, uuid)
	}, w, r)
}

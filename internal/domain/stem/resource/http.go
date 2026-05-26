package stems

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/stem"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Stems stem.Service
}

func New(stems stem.Service) *Resource {
	rc := Resource{
		Stems: stems,
	}

	routes := map[string]http.HandlerFunc{
		"GET /stems/{$}":           rc.List,
		"GET /stems/{uuid}":        rc.Get,
		"GET /stems/name/{name}":   rc.GetByName,
		"POST /stems/{$}":          rc.Create,
		"PUT /stems/rename/{uuid}": rc.Rename,
		"DELETE /stems/{uuid}":     rc.Delete,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Stems.List, w, r)
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (stem.Result, error) {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return stem.Result{}, resource.ErrBadUUID
		}

		return rc.Stems.Get(ctx, uuid)
	}, w, r)
}

func (rc *Resource) GetByName(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (stem.Result, error) {
		name := r.PathValue("name")
		return rc.Stems.GetByName(ctx, name)
	}, w, r)
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	resource.PostHandler(r.Context(), rc.Stems.Create, w, r)
}

func (rc *Resource) Rename(w http.ResponseWriter, r *http.Request) {
	resource.PutHandler(r.Context(), func(ctx context.Context, req stem.Rename) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Stems.Rename(ctx, uuid, req)
	}, w, r)
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	resource.DeleteHandler(r.Context(), func(ctx context.Context) error {
		uuid, err := uuid.FromString(r.PathValue("uuid"))
		if err != nil {
			return resource.ErrBadUUID
		}

		return rc.Stems.Delete(ctx, uuid)
	}, w, r)
}

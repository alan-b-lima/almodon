package stems

import (
	"context"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain/stem"
	"github.com/alan-b-lima/almodon/internal/support"
	"github.com/alan-b-lima/almodon/internal/support/resource"
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
		"GET /stems/{$}":         rc.Get,
		"GET /stems/history/{$}": rc.History,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), rc.Stems.Get, w, r)
}

func (rc *Resource) History(w http.ResponseWriter, r *http.Request) {
	resource.GetHandler(r.Context(), func(ctx context.Context) (struct{}, error) {
		return struct{}{}, support.ErrTODO
	}, w, r)
}

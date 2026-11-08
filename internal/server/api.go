package server

import (
	"fmt"
	"net/http"

	"github.com/alan-b-lima/almodon/internal/domain"
)

type API struct {
	http.ServeMux
	*domain.Domain
}

func NewAPI(opts ...domain.Option) (*API, error) {
	domain, err := domain.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("domain mounting: %w", err)
	}

	handlers := map[string]http.Handler{
		"auth":       domain.Resources.Auth,
		"items":      domain.Resources.Items,
		"materials":  domain.Resources.Materials,
		"promotions": domain.Resources.Promotions,
		"users":      domain.Resources.Users,
	}

	api := &API{Domain: domain}
	for name, handler := range handlers {
		api.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	return api, nil
}

func (a *API) Close() error {
	return a.Bundle.Close()
}

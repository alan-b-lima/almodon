package api

import (
	"net/http"

	"github.com/alan-b-lima/almodon/internal/auth"
	sessionrepo "github.com/alan-b-lima/almodon/internal/domain/session/repository"
	userrepo "github.com/alan-b-lima/almodon/internal/domain/user/repository"
	users "github.com/alan-b-lima/almodon/internal/domain/user/resource"
	userserve "github.com/alan-b-lima/almodon/internal/domain/user/service"
)

type router struct{ http.ServeMux }

func New() http.Handler {
	var r router

	var (
		repoSessions = sessionrepo.NewMap()
		repoUsers    = userrepo.NewMap()
	)

	serveUsers := userserve.New(userserve.NewService(repoUsers, repoSessions))

	resources := map[string]http.Handler{
		"users": users.New(serveUsers),
	}

	for name, handler := range resources {
		r.Handle("/api/v1/"+name+"/", http.StripPrefix("/api/v1", handler))
	}

	// temp
	{
		repoUsers.Create(1, "Alan Barbosa Lima", "alan-lima.al@ufvjm.edu.br", "12345678", auth.Chief)
		repoUsers.Create(2, "Breno Augusto Braga Oliveira", "b@ufvjm.edu.br", "12345678", auth.Admin)
		repoUsers.Create(3, "Lucas Rocha Oliveira", "l@ufvjm.edu.br", "12345678", auth.User)
		repoUsers.Create(4, "Luiz Felipe Melo Oliveira", "l@ufvjm.edu.br", "12345678", auth.Admin)
		repoUsers.Create(5, "Otávio Gomes Calazans", "o@ufvjm.edu.br", "12345678", auth.User)
		repoUsers.Create(6, "Rafael Gomes Silva", "r@ufvjm.edu.br", "12345678", auth.User)
	}

	return &r
}

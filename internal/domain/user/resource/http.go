package users

import (
	"net/http"
	"strconv"

	"github.com/alan-b-lima/almodon/internal/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/resource"
	"github.com/alan-b-lima/almodon/internal/xerrors"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Resource struct {
	http.ServeMux
	Users user.Service
}

func New(users user.Service) *Resource {
	rc := Resource{Users: users}

	routes := map[string]http.HandlerFunc{
		"GET /users/":              rc.List,
		"GET /users/{uuid}":        rc.Get,
		"GET /users/siape/{siape}": rc.GetBySIAPE,
		"POST /users/":             rc.Create,
		"PATCH /users/{uuid}":      rc.Patch,
		"DELETE /users/{uuid}":     rc.Delete,
		"POST /users/auth/":        rc.Authenticate,
		"GET /users/me/":           rc.Me,
	}

	for route, handler := range routes {
		rc.Handle(route, handler)
	}

	return &rc
}

func (rc *Resource) List(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	req := user.ListRequest{Offset: 0, Limit: 10}
	if err := resource.QueryParams(r.URL.Query(), &req); err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.List(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Get(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	uuid, err := uuid.FromString(r.PathValue("uuid"))
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}
	req := user.GetRequest{UUID: uuid}

	res, err := rc.Users.Get(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) GetBySIAPE(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	siape, err := strconv.Atoi(r.PathValue("siape"))
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}
	req := user.GetBySIAPERequest{SIAPE: siape}

	res, err := rc.Users.GetBySIAPE(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Create(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	var req user.CreateRequest

	if err := resource.DecodeJSON(&req, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Create(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Patch(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	uuid, err := uuid.FromString(r.PathValue("uuid"))
	if err != nil {
		resource.WriteJsonError(w, xerrors.ErrBadUUID)
		return
	}
	req := user.PatchRequest{UUID: uuid}

	if err := resource.DecodeJSON(&req, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Patch(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	uuid, err := uuid.FromString(r.PathValue("uuid"))
	if err != nil {
		resource.WriteJsonError(w, xerrors.ErrBadUUID)
		return
	}

	req := user.DeleteRequest{UUID: uuid}
	if err := rc.Users.Delete(act, req); err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rc *Resource) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req user.AuthRequest
	if err := resource.DecodeJSON(&req, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	res, err := rc.Users.Authenticate(req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     resource.SessionCookieName,
		Value:    res.UUID.String(),
		Expires:  res.Expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	if err := resource.EncodeJSON(&res, http.StatusCreated, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func (rc *Resource) Me(w http.ResponseWriter, r *http.Request) {
	act, err := resource.Session(rc.Users, r)
	if err != nil {
		resource.WriteJsonError(w, xerrors.ErrUnauthenticatedUser.New(err))
		return
	}

	req := user.GetRequest{UUID: act.User()}
	res, err := rc.Users.Get(act, req)
	if err != nil {
		resource.WriteJsonError(w, err)
		return
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		resource.WriteJsonError(w, err)
		return
	}
}

func n[Res, Req any](u user.Service, req Req, fn func(auth.Actor, Req) (Res, error), w http.ResponseWriter, r *http.Request) error {
	act, err := resource.Session(u, r)
	if err != nil {
		return err
	}

	if err := resource.DecodeJSON(&req, r); err != nil {
		return err
	}

	res, err := fn(act, req)
	if err != nil {
		return err
	}

	if err := resource.EncodeJSON(&res, http.StatusOK, w, r); err != nil {
		return err
	}

	return nil
}

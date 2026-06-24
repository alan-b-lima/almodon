package sessions

import (
	"context"
	"net/http"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/resource"
)

// on malformed session token, clear the token and let the user proceed as unlogged.
type Handler struct {
	Handler http.Handler
}

func Wrap(handler http.Handler) http.Handler {
	return &Handler{Handler: handler}
}

func WrapFunc(handler http.HandlerFunc) http.Handler {
	return Wrap(handler)
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, err := Session(ctx, r)
	if err != nil {
		if err != session.ErrInvalidToken {
			resource.WriteError(w, err)
			return
		}

		DeleteCookie(w)
	}

	r = r.WithContext(ctx)
	s.Handler.ServeHTTP(w, r)
}

const SessionIdentifier = "session"

func Session(ctx context.Context, r *http.Request) (context.Context, error) {
	token, err := Cookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			return ctx, nil
		}

		return nil, err
	}

	return context.WithValue(ctx, "session", token), nil
}

func Cookie(r *http.Request) (session.Token, error) {
	s, err := r.Cookie(SessionIdentifier)
	if err != nil {
		return session.Token{}, http.ErrNoCookie
	}

	token, err := session.ParseString(s.Value)
	if err != nil {
		return session.Token{}, session.ErrInvalidToken
	}

	return token, nil
}

func SetCookie(w http.ResponseWriter, token session.Token, expires time.Time) {
	cookie := &http.Cookie{
		Name:     SessionIdentifier,
		Value:    token.String(),
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     SessionIdentifier,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

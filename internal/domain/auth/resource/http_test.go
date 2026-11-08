package auths_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/server"
)

func init() {
	api, err := server.NewAPI(domain.InMemory, ^domain.RootUser)
	if err != nil {
		panic(err)
	}

	root_user := user.Create{
		SIAPE:    "0000000",
		Name:     "Raiz",
		Email:    "noreply@ufvjm.edu.br",
		Password: "12345678",
		Role:     auth.Maintainer,
	}

	_, err = api.Cores.Users.Create(context.Background(), root_user)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	url := url.URL{
		Scheme: "http",
		Host:   ln.Addr().String(),
	}

	Origin = url.String()
	go http.Serve(ln, api)

	time.Sleep(300 * time.Millisecond)
}

var Origin string

func TestLogin(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/auth/",
		strings.NewReader(`{"siape":"0000000","password":"12345678"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(resp.StatusCode, err)
		}

		t.Fatalf("status=%d body=%s", resp.StatusCode, body)
	}

	if len(resp.Cookies()) == 0 {
		t.Fatal("expected session cookie after login")
	}
}

func TestLoginWithInvalidPassword(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/auth/",
		strings.NewReader(`{"siape":"0000000","password":"senha-incorreta"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(resp.StatusCode, err)
		}

		t.Fatalf("expected authentication error, got status=%d body=%s", resp.StatusCode, body)
	}
}

func TestLogout(t *testing.T) {
	var client http.Client

	login_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/auth/",
		strings.NewReader(`{"siape":"0000000","password":"12345678"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	login_req.Header.Add("Content-Type", "application/json")
	login_req.Header.Add("Accept", "application/json")

	login_resp, err := client.Do(login_req)
	if err != nil {
		t.Fatal(err)
	}
	defer login_resp.Body.Close()

	if login_resp.StatusCode >= 400 {
		body, err := io.ReadAll(login_resp.Body)
		if err != nil {
			t.Fatal(login_resp.StatusCode, err)
		}

		t.Fatalf("status=%d body=%s", login_resp.StatusCode, body)
	}

	cookies := login_resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie after login")
	}

	logout_req, err := http.NewRequest(
		http.MethodDelete,
		Origin+"/api/v1/auth/",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	logout_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		logout_req.AddCookie(cookie)
	}

	logout_resp, err := client.Do(logout_req)
	if err != nil {
		t.Fatal(err)
	}
	defer logout_resp.Body.Close()

	if logout_resp.StatusCode >= 400 {
		body, err := io.ReadAll(logout_resp.Body)
		if err != nil {
			t.Fatal(logout_resp.StatusCode, err)
		}

		t.Fatalf("status=%d body=%s", logout_resp.StatusCode, body)
	}
}

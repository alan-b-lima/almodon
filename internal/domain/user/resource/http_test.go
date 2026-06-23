package users_test

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/alan-b-lima/almodon/internal/almodon"
	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/internal/domain/auth"
	sessions "github.com/alan-b-lima/almodon/internal/domain/session/resource"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

func init() {
	api, err := almodon.NewAPI(domain.InMemory, ^domain.RootUser)
	if err != nil {
		panic(err)
	}

	root_user := user.Create{
		SIAPE:    "0000000",
		Name:     "Raiz",
		Email:    "noreply@ufvjm.edu.br",
		Password: "12345678",
		Role:     auth.Chief,
	}

	_, err = api.Cores.Users.Create(context.Background(), root_user)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	url := url.URL{
		Scheme: "http",
		Host:   ln.Addr().String(),
	}

	Origin = url.String()
	go http.Serve(ln, sessions.Wrap(api))

	time.Sleep(300 * time.Millisecond)
}

var Origin string

func login(t *testing.T) []*http.Cookie {
	t.Helper()

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

	var client http.Client

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

		t.Fatalf("login: status=%d body=%s", login_resp.StatusCode, body)
	}

	cookies := login_resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("login: expected session cookie")
	}

	return cookies
}

func addCookies(req *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
}

func createUser(t *testing.T, cookies []*http.Cookie, siape, name, email string) user.CreateResult {
	t.Helper()

	body := `{"siape":"` + siape + `","name":"` + name + `","email":"` + email + `","password":"senha123","role":"user"}`

	create_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/users/",
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	create_req.Header.Add("Content-Type", "application/json")
	create_req.Header.Add("Accept", "application/json")

	addCookies(create_req, cookies)

	var client http.Client

	create_resp, err := client.Do(create_req)
	if err != nil {
		t.Fatal(err)
	}
	defer create_resp.Body.Close()

	if create_resp.StatusCode >= 400 {
		body, err := io.ReadAll(create_resp.Body)
		if err != nil {
			t.Fatal(create_resp.StatusCode, err)
		}

		t.Fatalf("create user: status=%d body=%s", create_resp.StatusCode, body)
	}

	var result user.CreateResult
	if err := json.NewDecoder(create_resp.Body).Decode(&result); err != nil {
		t.Fatalf("create user: decode response: %v", err)
	}

	if result.UUID == (uuid.UUID{}) {
		t.Fatal("create user: expected non-zero UUID in response")
	}

	return result
}

func getUserBySIAPE(t *testing.T, cookies []*http.Cookie, siape string) user.Result {
	t.Helper()

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/users/siape/"+siape,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	addCookies(get_req, cookies)

	var client http.Client

	get_resp, err := client.Do(get_req)
	if err != nil {
		t.Fatal(err)
	}
	defer get_resp.Body.Close()

	if get_resp.StatusCode >= 400 {
		body, err := io.ReadAll(get_resp.Body)
		if err != nil {
			t.Fatal(get_resp.StatusCode, err)
		}

		t.Fatalf("get user by siape: status=%d body=%s", get_resp.StatusCode, body)
	}

	var result user.Result
	if err := json.NewDecoder(get_resp.Body).Decode(&result); err != nil {
		t.Fatalf("get user by siape: decode response: %v", err)
	}

	return result
}

func getUserByUUID(t *testing.T, cookies []*http.Cookie, user_uuid string) user.Result {
	t.Helper()

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/users/"+user_uuid,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	addCookies(get_req, cookies)

	var client http.Client

	get_resp, err := client.Do(get_req)
	if err != nil {
		t.Fatal(err)
	}
	defer get_resp.Body.Close()

	if get_resp.StatusCode >= 400 {
		body, err := io.ReadAll(get_resp.Body)
		if err != nil {
			t.Fatal(get_resp.StatusCode, err)
		}

		t.Fatalf("get user by uuid: status=%d body=%s", get_resp.StatusCode, body)
	}

	var result user.Result
	if err := json.NewDecoder(get_resp.Body).Decode(&result); err != nil {
		t.Fatalf("get user by uuid: decode response: %v", err)
	}

	return result
}

func TestCreateUser(t *testing.T) {
	cookies := login(t)

	result := createUser(t, cookies, "1234567", "Bauru", "bauruzete@ufvjm.edu.br")

	if result.UUID == (uuid.UUID{}) {
		t.Fatal("create user: expected non-zero UUID in response")
	}
}

func TestGetUserBySIAPE(t *testing.T) {
	cookies := login(t)

	create_result := createUser(t, cookies, "2345678", "João Paulo", "joao.paulo@ufvjm.edu.br")
	result := getUserBySIAPE(t, cookies, "2345678")

	if result.UUID != create_result.UUID {
		t.Fatalf("get user by siape: expected uuid=%v, got %v", create_result.UUID, result.UUID)
	}

	if result.SIAPE != "2345678" {
		t.Fatalf("get user by siape: expected siape=%q, got %q", "2345678", result.SIAPE)
	}

	if result.Name != "João Paulo" {
		t.Fatalf("get user by siape: expected name=%q, got %q", "João Paulo", result.Name)
	}

	if result.Email != "joao.paulo@ufvjm.edu.br" {
		t.Fatalf("get user by siape: expected email=%q, got %q", "joao.paulo@ufvjm.edu.br", result.Email)
	}
}

func TestPatchUser(t *testing.T) {
	cookies := login(t)

	create_result := createUser(t, cookies, "3456789", "Alan Turing", "turing.alan@ufvjm.edu.br")

	patch_req, err := http.NewRequest(
		http.MethodPatch,
		Origin+"/api/v1/users/"+create_result.UUID.String(),
		strings.NewReader(`{"name":"Alan Editado","email":"turing.alan.novo@ufvjm.edu.br"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	patch_req.Header.Add("Content-Type", "application/json")
	patch_req.Header.Add("Accept", "application/json")

	addCookies(patch_req, cookies)

	var client http.Client

	patch_resp, err := client.Do(patch_req)
	if err != nil {
		t.Fatal(err)
	}
	defer patch_resp.Body.Close()

	if patch_resp.StatusCode >= 400 {
		body, err := io.ReadAll(patch_resp.Body)
		if err != nil {
			t.Fatal(patch_resp.StatusCode, err)
		}

		t.Fatalf("patch user: status=%d body=%s", patch_resp.StatusCode, body)
	}

	result := getUserBySIAPE(t, cookies, "3456789")

	if result.UUID != create_result.UUID {
		t.Fatalf("patch user: expected uuid=%v, got %v", create_result.UUID, result.UUID)
	}

	if result.Name != "Alan Editado" {
		t.Fatalf("patch user: expected name=%q, got %q", "Alan Editado", result.Name)
	}

	if result.Email != "turing.alan.novo@ufvjm.edu.br" {
		t.Fatalf("patch user: expected email=%q, got %q", "turing.alan.novo@ufvjm.edu.br", result.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	cookies := login(t)

	create_result := createUser(t, cookies, "4567890", "Luan Deletado", "deleted.luan@ufvjm.edu.br")

	delete_req, err := http.NewRequest(
		http.MethodDelete,
		Origin+"/api/v1/users/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	delete_req.Header.Add("Accept", "application/json")

	addCookies(delete_req, cookies)

	var client http.Client

	delete_resp, err := client.Do(delete_req)
	if err != nil {
		t.Fatal(err)
	}
	defer delete_resp.Body.Close()

	if delete_resp.StatusCode >= 400 {
		body, err := io.ReadAll(delete_resp.Body)
		if err != nil {
			t.Fatal(delete_resp.StatusCode, err)
		}

		t.Fatalf("delete user: status=%d body=%s", delete_resp.StatusCode, body)
	}

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/users/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	addCookies(get_req, cookies)

	get_resp, err := client.Do(get_req)
	if err != nil {
		t.Fatal(err)
	}
	defer get_resp.Body.Close()

	if get_resp.StatusCode < 400 {
		body, err := io.ReadAll(get_resp.Body)
		if err != nil {
			t.Fatal(get_resp.StatusCode, err)
		}

		t.Fatalf("get deleted user: expected error, got status=%d body=%s", get_resp.StatusCode, body)
	}
}

func TestGetUserByUUID(t *testing.T) {
	cookies := login(t)

	create_result := createUser(t, cookies, "6789012", "Vitor o Grandioso UUID", "vitor.grande.uuid@ufvjm.edu.br")
	result := getUserByUUID(t, cookies, create_result.UUID.String())

	if result.UUID != create_result.UUID {
		t.Fatalf("get user by uuid: expected uuid=%v, got %v", create_result.UUID, result.UUID)
	}

	if result.SIAPE != "6789012" {
		t.Fatalf("get user by uuid: expected siape=%q, got %q", "6789012", result.SIAPE)
	}

	if result.Name != "Vitor o Grandioso UUID" {
		t.Fatalf("get user by uuid: expected name=%q, got %q", "Vitor o Grandioso UUID", result.Name)
	}

	if result.Email != "vitor.grande.uuid@ufvjm.edu.br" {
		t.Fatalf("get user by uuid: expected email=%q, got %q", "vitor.grande.uuid@ufvjm.edu.br", result.Email)
	}
}

func TestCreateUserWithoutSession(t *testing.T) {
	create_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/users/",
		strings.NewReader(`{"siape":"5678901","name":"Bloqueado","email":"bloqueado@ufvjm.edu.br","password":"senha123","role":"user"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	create_req.Header.Add("Content-Type", "application/json")
	create_req.Header.Add("Accept", "application/json")

	var client http.Client

	create_resp, err := client.Do(create_req)
	if err != nil {
		t.Fatal(err)
	}
	defer create_resp.Body.Close()

	if create_resp.StatusCode < 400 {
		body, err := io.ReadAll(create_resp.Body)
		if err != nil {
			t.Fatal(create_resp.StatusCode, err)
		}

		t.Fatalf("expected auth error without session, got status=%d body=%s", create_resp.StatusCode, body)
	}
}

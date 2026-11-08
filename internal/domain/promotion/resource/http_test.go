package promotions_test

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

	"github.com/alan-b-lima/almodon/internal/domain"
	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	sessions "github.com/alan-b-lima/almodon/internal/domain/session/resource"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	server "github.com/alan-b-lima/almodon/internal/server"
	"github.com/alan-b-lima/almodon/pkg/uuid"
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

func createUser(t *testing.T, cookies []*http.Cookie, siape, name, email string) string {
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

	for _, cookie := range cookies {
		create_req.AddCookie(cookie)
	}

	var client http.Client

	create_resp, err := client.Do(create_req)
	if err != nil {
		t.Fatal(err)
	}
	defer create_resp.Body.Close()

	if create_resp.StatusCode >= 400 {
		b, err := io.ReadAll(create_resp.Body)
		if err != nil {
			t.Fatal(create_resp.StatusCode, err)
		}

		t.Fatalf("create user: status=%d body=%s", create_resp.StatusCode, b)
	}

	var result user.CreateResult
	if err := json.NewDecoder(create_resp.Body).Decode(&result); err != nil {
		t.Fatalf("create user: decode response: %v", err)
	}

	return result.UUID.String()
}

func createPromotion(t *testing.T, cookies []*http.Cookie, user_uuid string) promotion.CreateResult {
	t.Helper()

	create_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/promotions/",
		strings.NewReader(`{"user":"`+user_uuid+`"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	create_req.Header.Add("Content-Type", "application/json")
	create_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		create_req.AddCookie(cookie)
	}

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

		t.Fatalf("create promotion: status=%d body=%s", create_resp.StatusCode, body)
	}

	var result promotion.CreateResult
	if err := json.NewDecoder(create_resp.Body).Decode(&result); err != nil {
		t.Fatalf("create promotion: decode response: %v", err)
	}

	if result.UUID == (uuid.UUID{}) {
		t.Fatal("create promotion: expected non-zero UUID in response")
	}

	return result
}

func TestCreatePromotion(t *testing.T) {
	cookies := login(t)

	user_uuid := createUser(t, cookies, "1111111", "Alanzão", "alan.promotion@ufvjm.edu.br")

	create_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/promotions/",
		strings.NewReader(`{"user":"`+user_uuid+`"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	create_req.Header.Add("Content-Type", "application/json")
	create_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		create_req.AddCookie(cookie)
	}

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

		t.Fatalf("create promotion: status=%d body=%s", create_resp.StatusCode, body)
	}

	var result promotion.CreateResult
	if err := json.NewDecoder(create_resp.Body).Decode(&result); err != nil {
		t.Fatalf("create promotion: decode response: %v", err)
	}

	if result.UUID == (uuid.UUID{}) {
		t.Fatal("create promotion: expected non-zero UUID in response")
	}
}

func TestGetPromotion(t *testing.T) {
	cookies := login(t)

	user_uuid := createUser(t, cookies, "2222222", "Luan", "luan.promotion@ufvjm.edu.br")
	create_result := createPromotion(t, cookies, user_uuid)

	var client http.Client

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/promotions/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		get_req.AddCookie(cookie)
	}

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

		t.Fatalf("get promotion: status=%d body=%s", get_resp.StatusCode, body)
	}

	var result promotion.Result
	if err := json.NewDecoder(get_resp.Body).Decode(&result); err != nil {
		t.Fatalf("get promotion: decode response: %v", err)
	}

	if result.UUID != create_result.UUID {
		t.Fatalf("get promotion: expected uuid=%v, got %v", create_result.UUID, result.UUID)
	}

	if result.User.String() != user_uuid {
		t.Fatalf("get promotion: expected user=%v, got %v", user_uuid, result.User)
	}

	if result.Expires.IsZero() {
		t.Fatal("get promotion: expected non-zero expires")
	}
}

func TestUpdatePromotion(t *testing.T) {
	cookies := login(t)

	user_uuid := createUser(t, cookies, "3333333", "Bauru", "bauru.promotion@ufvjm.edu.br")
	create_result := createPromotion(t, cookies, user_uuid)

	var client http.Client

	update_req, err := http.NewRequest(
		http.MethodPut,
		Origin+"/api/v1/promotions/"+create_result.UUID.String(),
		strings.NewReader(`{"max_age":7200000000000}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	update_req.Header.Add("Content-Type", "application/json")
	update_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		update_req.AddCookie(cookie)
	}

	update_resp, err := client.Do(update_req)
	if err != nil {
		t.Fatal(err)
	}
	defer update_resp.Body.Close()

	if update_resp.StatusCode >= 400 {
		body, err := io.ReadAll(update_resp.Body)
		if err != nil {
			t.Fatal(update_resp.StatusCode, err)
		}

		t.Fatalf("update promotion: status=%d body=%s", update_resp.StatusCode, body)
	}

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/promotions/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		get_req.AddCookie(cookie)
	}

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

		t.Fatalf("get updated promotion: status=%d body=%s", get_resp.StatusCode, body)
	}

	var result promotion.Result
	if err := json.NewDecoder(get_resp.Body).Decode(&result); err != nil {
		t.Fatalf("get updated promotion: decode response: %v", err)
	}

	if result.UUID != create_result.UUID {
		t.Fatalf("get updated promotion: expected uuid=%v, got %v", create_result.UUID, result.UUID)
	}

	if result.Expires.IsZero() {
		t.Fatal("get updated promotion: expected non-zero expires")
	}
}

func TestDeletePromotion(t *testing.T) {
	cookies := login(t)

	user_uuid := createUser(t, cookies, "4444444", "JuanP", "juan.pablomotion@ufvjm.edu.br")
	create_result := createPromotion(t, cookies, user_uuid)

	var client http.Client

	delete_req, err := http.NewRequest(
		http.MethodDelete,
		Origin+"/api/v1/promotions/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	delete_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		delete_req.AddCookie(cookie)
	}

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

		t.Fatalf("delete promotion: status=%d body=%s", delete_resp.StatusCode, body)
	}

	get_req, err := http.NewRequest(
		http.MethodGet,
		Origin+"/api/v1/promotions/"+create_result.UUID.String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	get_req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		get_req.AddCookie(cookie)
	}

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

		t.Fatalf("get deleted promotion: expected error, got status=%d body=%s", get_resp.StatusCode, body)
	}
}

func TestCreatePromotionWithoutSession(t *testing.T) {
	cookies := login(t)

	user_uuid := createUser(t, cookies, "5555555", "Vitor O Grande", "vitor.promotion@ufvjm.edu.br")

	create_req, err := http.NewRequest(
		http.MethodPost,
		Origin+"/api/v1/promotions/",
		strings.NewReader(`{"user":"`+user_uuid+`"}`),
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

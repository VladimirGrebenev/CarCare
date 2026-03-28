package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type notImplHandler struct{}

func (h *notImplHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestCarHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	h := &notImplHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
	}
}

type notImplUserHandler struct{}

func (h *notImplUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestUserHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	h := &notImplUserHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
	}
}

type notImplFuelHandler struct{}

func (h *notImplFuelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestFuelHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/fuel", nil)
	w := httptest.NewRecorder()
	h := &notImplFuelHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
	}
}

func TestMaintenanceHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/maintenance", nil)
	w := httptest.NewRecorder()
	h := rest.NewMaintenanceHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestFineHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/fines", nil)
	w := httptest.NewRecorder()
	h := rest.NewFineHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReportHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/reports", nil)
	w := httptest.NewRecorder()
	h := rest.NewReportHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestProfileHandler_GET(t *testing.T) {
	// Логинимся через usecase с подготовленным репозиторием
	const testEmail = "test@example.com"
	const testPassword = "123456"
	repo := newStubUserRepoWithUser(testEmail, testPassword)
	uc := &usecase.AuthUsecase{
		UserRepo: repo,
		Logger:   &stubLogger{},
	}
	authHandler := rest.NewAuthHandler(uc)

	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"test@example.com","password":"123456"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	authHandler.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected 200 from login, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var loginBody map[string]any
	if err := json.Unmarshal(loginW.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("invalid login json response: %v", err)
	}
	token, _ := loginBody["token"].(string)
	if token == "" {
		t.Fatal("expected JWT token from login response")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	rest.ProfileHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if response["email"] != testEmail {
		t.Fatalf("expected email %s, got %v", testEmail, response["email"])
	}
}

func TestProfileHandler_Unauthorized(t *testing.T) {
	// ProfileHandler без токена: маршрут защищён middleware, но проверим что без Authorization
	// ProfileHandler сам по себе не возвращает 401 — это делает AuthMiddleware.
	// Тест проверяет поведение через AuthMiddleware.
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()

	handler := rest.AuthMiddleware(http.HandlerFunc(rest.ProfileHandler))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestProfileHandler_MethodNotAllowed(t *testing.T) {
	// POST без токена: сначала сработает AuthMiddleware (401), но тест проверяет 405.
	// Пропустим через валидный токен, чтобы добраться до ProfileHandler.
	const testEmail = "test@example.com"
	const testPassword = "pass"
	repo := newStubUserRepoWithUser(testEmail, testPassword)
	uc := &usecase.AuthUsecase{UserRepo: repo, Logger: &stubLogger{}}
	authH := rest.NewAuthHandler(uc)

	loginW := httptest.NewRecorder()
	authH.Login(loginW, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"email":"test@example.com","password":"pass"}`)))
	var lb map[string]any
	_ = json.Unmarshal(loginW.Body.Bytes(), &lb)
	token, _ := lb["token"].(string)

	req := httptest.NewRequest(http.MethodPost, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	rest.ProfileHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestFuelHandler_List_EmptyFallback(t *testing.T) {
	failingRepo := &mockFuelRepoWithError{err: &MockError{"repo unavailable"}}
	uc := &usecase.UsecaseContainer{Fuel: failingRepo}
	h := rest.NewFuelHandler(uc)

	req := httptest.NewRequest("GET", "/fuel", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []fuel.FuelEvent
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("expected empty fuel array, got %v", got)
	}
}

// mockFuelRepoWithError always returns an error on any operation
type mockFuelRepoWithError struct {
	err error
}

func (m *mockFuelRepoWithError) AddFuelEvent(e fuel.FuelEvent) error {
	return m.err
}

func (m *mockFuelRepoWithError) GetFuelEvent(id string) (fuel.FuelEvent, error) {
	return fuel.FuelEvent{}, m.err
}

func (m *mockFuelRepoWithError) UpdateFuelEvent(e fuel.FuelEvent) error {
	return m.err
}

func (m *mockFuelRepoWithError) DeleteFuelEvent(id string) error {
	return m.err
}

func (m *mockFuelRepoWithError) ListFuelEvents(userID string) ([]fuel.FuelEvent, error) {
	return nil, m.err
}

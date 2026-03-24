package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func TestRegisterSuccess(t *testing.T) {
	uc := &usecase.AuthUsecase{
		UserRepo:    &stubUserRepo{},
		EmailSender: &stubEmailSender{},
		Logger:      &stubLogger{},
	}
	h := rest.NewAuthHandler(uc)
	r := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"email":"a@b.com","password":"123"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if _, ok := body["token"]; !ok {
		t.Errorf("expected token field in response")
	}
	if _, ok := body["user"]; !ok {
		t.Errorf("expected user field in response")
	}
}

func TestRegisterMethodNotAllowed(t *testing.T) {
	uc := &usecase.AuthUsecase{Logger: &stubLogger{}}
	h := rest.NewAuthHandler(uc)
	r := httptest.NewRequest(http.MethodGet, "/api/auth/register", nil)
	w := httptest.NewRecorder()

	h.Register(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}
}

func TestLoginSuccess(t *testing.T) {
	uc := &usecase.AuthUsecase{Logger: &stubLogger{}}
	h := rest.NewAuthHandler(uc)
	r := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"a@b.com","password":"123"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if _, ok := body["token"]; !ok {
		t.Errorf("expected token field in response")
	}
	if _, ok := body["user"]; !ok {
		t.Errorf("expected user field in response")
	}
}

func TestLoginMethodNotAllowed(t *testing.T) {
	uc := &usecase.AuthUsecase{Logger: &stubLogger{}}
	h := rest.NewAuthHandler(uc)
	r := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}
}

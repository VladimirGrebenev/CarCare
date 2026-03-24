package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func TestAliasPrefixHandler_RewritesPath(t *testing.T) {
	h := rest.AliasPrefixHandler("/api", "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cars/123" {
			t.Fatalf("expected rewritten path /cars/123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/cars/123", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestOAuthProvider_Google(t *testing.T) {
	h := rest.NewAuthHandler(&usecase.AuthUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/oauth/google", nil)
	w := httptest.NewRecorder()

	h.OAuthProvider(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json, got err: %v", err)
	}
	if payload["error"] == "" {
		t.Fatalf("expected error payload, got: %v", payload)
	}
}

func TestOAuthProvider_UnknownProvider(t *testing.T) {
	h := rest.NewAuthHandler(&usecase.AuthUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/oauth/unknown", nil)
	w := httptest.NewRecorder()

	h.OAuthProvider(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestOAuthProvider_MethodNotAllowed(t *testing.T) {
	h := rest.NewAuthHandler(&usecase.AuthUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/oauth/google", nil)
	w := httptest.NewRecorder()

	h.OAuthProvider(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

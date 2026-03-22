package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func TestRegisterHandler(t *testing.T) {
	h := rest.NewAuthHandler(&usecase.AuthUsecase{})
	r := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{"email":"a@b.com","password":"123"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Errorf("expected 201 or 200, got %d", resp.StatusCode)
	}
}

// ...интеграционные тесты для ConfirmEmail, Login, ForgotPassword, ResetPassword, OAuth, Refresh, Logout

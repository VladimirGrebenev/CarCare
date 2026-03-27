package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// AuthHandler реализует auth-flow: регистрация, подтверждение email, восстановление пароля, OAuth, JWT/сессии, rate limiting
// TODO: реализовать централизованное логирование

type AuthHandler struct {
	UC     *usecase.AuthUsecase
	Logger usecase.Logger
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{UC: uc, Logger: uc.Logger}
}


func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	if h.Logger != nil {
		h.Logger.Info("Register endpoint")
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "invalid json")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeBadRequest(w, "email and password are required")
		return
	}

	token, err := h.UC.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user already exists"})
			return
		}
		writeServerError(w)
		return
	}

	writeAuthResponse(w, http.StatusCreated, token, req.Email)
}

func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ConfirmEmail endpoint")
	// TODO: получить token из запроса, вызвать usecase
}

func (h *AuthHandler) ResendConfirmation(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ResendConfirmation endpoint")
	// TODO: получить email, вызвать usecase
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w)
		return
	}

	if h.Logger != nil {
		h.Logger.Info("Login endpoint")
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "invalid json")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeBadRequest(w, "email and password are required")
		return
	}

	token, err := h.UC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	writeAuthResponse(w, http.StatusOK, token, req.Email)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ForgotPassword endpoint")
	// TODO: получить email, вызвать usecase
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ResetPassword endpoint")
	// TODO: получить token/newPassword, вызвать usecase
}

func (h *AuthHandler) OAuthYandex(w http.ResponseWriter, r *http.Request) {
	if h.Logger != nil {
		h.Logger.Info("OAuthYandex endpoint")
	}
	writeNotImplementedJSON(w, "oauth yandex not implemented")
}

func (h *AuthHandler) OAuthGoogle(w http.ResponseWriter, r *http.Request) {
	if h.Logger != nil {
		h.Logger.Info("OAuthGoogle endpoint")
	}
	writeNotImplementedJSON(w, "oauth google not implemented")
}

func (h *AuthHandler) OAuthProvider(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	provider := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	switch provider {
	case "google":
		h.OAuthGoogle(w, r)
	case "yandex":
		h.OAuthYandex(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unknown oauth provider"})
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("RefreshToken endpoint")
	// TODO: получить refreshToken, вызвать usecase
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Logout endpoint")
	// TODO: получить token, вызвать usecase
}

type authResponse struct {
	Token string            `json:"token"`
	User  map[string]string `json:"user"`
}

func writeAuthResponse(w http.ResponseWriter, status int, token, email string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if email == "" {
		email = "test@example.com"
	}
	_ = json.NewEncoder(w).Encode(authResponse{
		Token: token,
		User: map[string]string{
			"id":    "mock-user-id",
			"email": email,
			"name":  "Test User",
		},
	})
}

func writeBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow", http.MethodPost)
	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

func writeNotImplementedJSON(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

var rateLimiters = make(map[string]*rateLimiter)
var rlMu sync.Mutex

type rateLimiter struct {
	last  time.Time
	count int
}

// RateLimit — простая in-memory реализация (на IP+endpoint, 5 req/min)
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr + r.URL.Path
		rlMu.Lock()
		rl, ok := rateLimiters[key]
		if !ok || time.Since(rl.last) > time.Minute {
			rl = &rateLimiter{last: time.Now(), count: 1}
			rateLimiters[key] = rl
		} else {
			rl.count++
		}
		rl.last = time.Now()
		c := rl.count
		rlMu.Unlock()
		if c > 5 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

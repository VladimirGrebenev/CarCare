package rest

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// AuthHandler реализует auth-flow: регистрация, подтверждение email, восстановление пароля, OAuth, JWT/сессии, rate limiting
// TODO: реализовать централизованное логирование

type AuthHandler struct {
	UC *usecase.AuthUsecase
	Logger usecase.Logger
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{UC: uc, Logger: uc.Logger}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Register endpoint")
	// ...
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
	h.Logger.Info("Login endpoint")
	// ...
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
	h.Logger.Info("OAuthYandex endpoint")
	// ...
}

func (h *AuthHandler) OAuthGoogle(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("OAuthGoogle endpoint")
	// ...
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("RefreshToken endpoint")
	// TODO: получить refreshToken, вызвать usecase
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Logout endpoint")
	// TODO: получить token, вызвать usecase
}

// Rate limiting middleware
import (
       "sync"
       "time"
)

var rateLimiters = make(map[string]*rateLimiter)
var rlMu sync.Mutex

type rateLimiter struct {
       last time.Time
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

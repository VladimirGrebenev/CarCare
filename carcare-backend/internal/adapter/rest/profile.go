package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ProfileHandler возвращает профиль текущего пользователя.
// Маршрут защищён AuthMiddleware, поэтому здесь только читаем claims из токена.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	email, userID := extractClaimsFromRequest(r)

	profile := map[string]interface{}{
		"id":    userID,
		"name":  email,
		"email": email,
		"role":  "user",
		"cars":  []map[string]interface{}{},
	}
	_ = json.NewEncoder(w).Encode(profile)
}

// extractClaimsFromRequest разбирает JWT из заголовка Authorization и возвращает email и sub.
// Ошибки игнорируются — маршрут уже прошёл через AuthMiddleware.
func extractClaimsFromRequest(r *http.Request) (email, sub string) {
	header := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(getJWTSecret()), nil
	})
	if err != nil || !token.Valid {
		return "", ""
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if e, ok := claims["email"].(string); ok {
			email = e
		}
		if s, ok := claims["sub"].(string); ok {
			sub = s
		}
	}
	return email, sub
}

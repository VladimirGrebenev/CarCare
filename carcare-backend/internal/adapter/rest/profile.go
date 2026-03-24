package rest

import (
	"encoding/json"
	"net/http"
)

// ProfileHandler returns a test user profile as JSON
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	email, ok := validateBearerToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	profile := map[string]interface{}{
		"id":    1,
		"name":  "Test User",
		"email": email,
		"role":  "user",
		"cars":  []map[string]interface{}{},
	}
	_ = json.NewEncoder(w).Encode(profile)
}

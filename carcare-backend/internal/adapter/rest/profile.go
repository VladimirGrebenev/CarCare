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
	profile := map[string]interface{}{
		"id":    1,
		"name":  "Test User",
		"email": "test@example.com",
		"role":  "user",
	}
	json.NewEncoder(w).Encode(profile)
}

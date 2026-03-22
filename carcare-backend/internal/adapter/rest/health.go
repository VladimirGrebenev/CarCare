package rest

import (
	"net/http"
)

// HealthCheckHandler provides a simple health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TODO: Add REST handlers for car, fuel, fine, maintenance, user

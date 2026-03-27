package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// newMockUsecaseContainer returns a UsecaseContainer backed by in-memory mock repos.
func newMockUsecaseContainer() *usecase.UsecaseContainer {
	return &usecase.UsecaseContainer{
		Car:         NewMockCarRepo(),
		Fuel:        &mockFuelRepo{},
		Fine:        &mockFineRepo{},
		Maintenance: &mockMaintenanceRepo{},
	}
}

// TestMaintenanceHandler_GET tests the maintenance endpoint
func TestMaintenanceHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/maintenance", nil)
	w := httptest.NewRecorder()

	h := rest.NewMaintenanceHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}

	var body []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

// TestMaintenanceHandler_MethodNotAllowed tests non-supported methods
func TestMaintenanceHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/maintenance", nil)
	w := httptest.NewRecorder()

	h := rest.NewMaintenanceHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// TestFineHandler_GET tests the fines endpoint
func TestFineHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fines", nil)
	w := httptest.NewRecorder()

	h := rest.NewFineHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}

	var body []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

// TestFineHandler_MethodNotAllowed tests non-supported methods
func TestFineHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/fines", nil)
	w := httptest.NewRecorder()

	h := rest.NewFineHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// TestReportHandler_GET tests the reports endpoint
func TestReportHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	w := httptest.NewRecorder()

	h := rest.NewReportHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected content-type application/json, got %s", got)
	}
}

// TestReportHandler_MethodNotAllowed tests non-GET methods
func TestReportHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/reports", nil)
	w := httptest.NewRecorder()

	h := rest.NewReportHandler(newMockUsecaseContainer())
	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

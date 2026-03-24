package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
)

// TestMaintenanceHandler_GET tests the maintenance endpoint
func TestMaintenanceHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/maintenance", nil)
	w := httptest.NewRecorder()

	rest.MaintenanceHandler(w, req)
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
	if len(body) != 0 {
		t.Errorf("expected empty list, got %d items", len(body))
	}
}

// TestMaintenanceHandler_MethodNotAllowed tests non-GET methods
func TestMaintenanceHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/maintenance", nil)
	w := httptest.NewRecorder()

	rest.MaintenanceHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// TestFineHandler_GET tests the fines endpoint
func TestFineHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/fines", nil)
	w := httptest.NewRecorder()

	rest.FineHandler(w, req)
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
	if len(body) != 0 {
		t.Errorf("expected empty list, got %d items", len(body))
	}
}

// TestFineHandler_MethodNotAllowed tests non-GET methods
func TestFineHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/fines", nil)
	w := httptest.NewRecorder()

	rest.FineHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// TestReportHandler_GET tests the reports endpoint
func TestReportHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/reports", nil)
	w := httptest.NewRecorder()

	rest.ReportHandler(w, req)
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
	if len(body) != 0 {
		t.Errorf("expected empty list, got %d items", len(body))
	}
}

// TestReportHandler_MethodNotAllowed tests non-GET methods
func TestReportHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/reports", nil)
	w := httptest.NewRecorder()

	rest.ReportHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

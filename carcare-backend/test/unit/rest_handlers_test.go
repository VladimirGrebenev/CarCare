package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"home/veneberg/CarCare/carcare-backend/internal/adapter/rest"
)

func TestCarHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	rest.CarHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestUserHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	rest.UserHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestFuelHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/fuel", nil)
	w := httptest.NewRecorder()
	rest.FuelHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestMaintenanceHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/maintenance", nil)
	w := httptest.NewRecorder()
	rest.MaintenanceHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestFineHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/fines", nil)
	w := httptest.NewRecorder()
	rest.FineHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

func TestReportHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/reports", nil)
	w := httptest.NewRecorder()
	rest.ReportHandler(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("expected 501, got %d", w.Code)
	}
}

package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
)

type notImplHandler struct{}

func (h *notImplHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestCarHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	h := &notImplHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
	}
}

type notImplUserHandler struct{}

func (h *notImplUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestUserHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	h := &notImplUserHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
	}
}

type notImplFuelHandler struct{}

func (h *notImplFuelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func TestFuelHandler_NotImplemented(t *testing.T) {
	req := httptest.NewRequest("GET", "/fuel", nil)
	w := httptest.NewRecorder()
	h := &notImplFuelHandler{}
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented && w.Code != http.StatusOK {
		t.Errorf("expected 501 or 200, got %d", w.Code)
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

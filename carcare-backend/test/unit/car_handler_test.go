package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func setupCarHandler() *rest.CarHandler {
	repo := NewMockCarRepo()
	uc := &usecase.UsecaseContainer{Car: repo}
	return rest.NewCarHandler(uc)
}

func TestCarHandler_Update(t *testing.T) {
	h := setupCarHandler()
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1"})
	upd := car.Car{ID: "1", Brand: "B", Model: "X", Year: 2022, VIN: "VIN1"}
	b, _ := json.Marshal(upd)
	req := httptest.NewRequest("PUT", "/cars", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got car.Car
	json.NewDecoder(w.Body).Decode(&got)
	if got.Brand != "B" || got.Year != 2022 {
		t.Errorf("update failed: %+v", got)
	}
}

func TestCarHandler_Delete(t *testing.T) {
	h := setupCarHandler()
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1"})
	req := httptest.NewRequest("DELETE", "/cars/1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	// Проверка, что удалено
	req = httptest.NewRequest("GET", "/cars/1", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", w.Code)
	}
}

func TestCarHandler_BadRequest(t *testing.T) {
	h := setupCarHandler()
	req := httptest.NewRequest("POST", "/cars", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

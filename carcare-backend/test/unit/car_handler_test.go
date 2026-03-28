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
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1", UserID: "user-1"})
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
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1", UserID: "user-1"})
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

func TestCarHandler_List_Success(t *testing.T) {
	h := setupCarHandler()
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1", UserID: "user-1"})
	_ = h.Add.Execute(car.Car{ID: "2", Brand: "B", Model: "N", Year: 2021, VIN: "VIN2", UserID: "user-1"})
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []car.Car
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 cars, got %d", len(got))
	}
}

func TestCarHandler_List_EmptyFallback(t *testing.T) {
	// Create a handler with a failing repository
	failingRepo := &mockCarRepoWithError{}
	uc := &usecase.UsecaseContainer{Car: failingRepo}
	h := rest.NewCarHandler(uc)
	
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	
	// Should return 200 with empty array, not 500
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []car.Car
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("expected empty car array, got %v", got)
	}
}

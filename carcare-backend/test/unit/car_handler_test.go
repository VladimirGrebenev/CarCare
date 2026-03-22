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

type mockCarRepo struct {
	cars map[string]car.Car
}

func newMockCarRepo() *mockCarRepo {
	return &mockCarRepo{cars: make(map[string]car.Car)}
}

func (m *mockCarRepo) AddCar(c car.Car) error {
	if _, exists := m.cars[c.ID]; exists {
		return ErrDuplicateID
	}
	m.cars[c.ID] = c
	return nil
}

func (m *mockCarRepo) GetCar(id string) (car.Car, error) {
	c, ok := m.cars[id]
	if !ok {
		return car.Car{}, ErrNotFound
	}
	return c, nil
}

func (m *mockCarRepo) UpdateCar(c car.Car) error {
	if _, ok := m.cars[c.ID]; !ok {
		return ErrNotFound
	}
	m.cars[c.ID] = c
	return nil
}

func (m *mockCarRepo) DeleteCar(id string) error {
	if _, ok := m.cars[id]; !ok {
		return ErrNotFound
	}
	delete(m.cars, id)
	return nil
}

func (m *mockCarRepo) ListCars() ([]car.Car, error) {
	out := make([]car.Car, 0, len(m.cars))
	for _, c := range m.cars {
		out = append(out, c)
	}
	return out, nil
}

var (
	ErrNotFound    = &mockError{"not found"}
	ErrDuplicateID = &mockError{"duplicate id"}
)

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }

func setupCarHandler() *rest.CarHandler {
	repo := newMockCarRepo()
	uc := &usecase.UsecaseContainer{Car: repo}
	return rest.NewCarHandler(uc)
}

func TestCarHandler_AddGet(t *testing.T) {
	h := setupCarHandler()
	carObj := car.Car{ID: "1", Brand: "Test", Model: "M", Year: 2020, VIN: "VIN1"}
	b, _ := json.Marshal(carObj)
	req := httptest.NewRequest("POST", "/cars", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	// Get
	req = httptest.NewRequest("GET", "/cars/1", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got car.Car
	json.NewDecoder(w.Body).Decode(&got)
	if got.ID != "1" || got.VIN != "VIN1" {
		t.Errorf("unexpected car: %+v", got)
	}
}

func TestCarHandler_List(t *testing.T) {
	h := setupCarHandler()
	_ = h.Add.Execute(car.Car{ID: "1", Brand: "A", Model: "M", Year: 2020, VIN: "VIN1"})
	_ = h.Add.Execute(car.Car{ID: "2", Brand: "B", Model: "N", Year: 2021, VIN: "VIN2"})
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var cars []car.Car
	json.NewDecoder(w.Body).Decode(&cars)
	if len(cars) != 2 {
		t.Errorf("expected 2 cars, got %d", len(cars))
	}
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

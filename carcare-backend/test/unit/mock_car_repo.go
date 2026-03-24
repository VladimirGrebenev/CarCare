package unit

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"

type MockCarRepo struct {
    Cars map[string]car.Car
}

func NewMockCarRepo() *MockCarRepo {
    return &MockCarRepo{Cars: make(map[string]car.Car)}
}

func (m *MockCarRepo) AddCar(c car.Car) error {
    if _, exists := m.Cars[c.ID]; exists {
        return ErrDuplicateID
    }
    m.Cars[c.ID] = c
    return nil
}

func (m *MockCarRepo) GetCar(id string) (car.Car, error) {
    c, ok := m.Cars[id]
    if !ok {
        return car.Car{}, ErrNotFound
    }
    return c, nil
}

func (m *MockCarRepo) UpdateCar(c car.Car) error {
    if _, ok := m.Cars[c.ID]; !ok {
        return ErrNotFound
    }
    m.Cars[c.ID] = c
    return nil
}

func (m *MockCarRepo) DeleteCar(id string) error {
    if _, ok := m.Cars[id]; !ok {
        return ErrNotFound
    }
    delete(m.Cars, id)
    return nil
}

func (m *MockCarRepo) ListCars() ([]car.Car, error) {
    out := make([]car.Car, 0, len(m.Cars))
    for _, c := range m.Cars {
        out = append(out, c)
    }
    return out, nil
}

var (
    ErrNotFound    = &MockError{"not found"}
    ErrDuplicateID = &MockError{"duplicate id"}
)

type MockError struct{ msg string }

func (e *MockError) Error() string { return e.msg }

// mockCarRepoWithError always returns an error on any operation
type mockCarRepoWithError struct{}

func (m *mockCarRepoWithError) AddCar(c car.Car) error {
	return &MockError{"repo unavailable"}
}

func (m *mockCarRepoWithError) GetCar(id string) (car.Car, error) {
	return car.Car{}, &MockError{"repo unavailable"}
}

func (m *mockCarRepoWithError) UpdateCar(c car.Car) error {
	return &MockError{"repo unavailable"}
}

func (m *mockCarRepoWithError) DeleteCar(id string) error {
	return &MockError{"repo unavailable"}
}

func (m *mockCarRepoWithError) ListCars() ([]car.Car, error) {
	return nil, &MockError{"repo unavailable"}
}
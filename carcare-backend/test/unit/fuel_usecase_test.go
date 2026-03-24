package unit

import (
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type mockFuelRepo struct{}

func (m *mockFuelRepo) AddFuelEvent(e fuel.FuelEvent) error { return nil }
func (m *mockFuelRepo) GetFuelEvent(id string) (fuel.FuelEvent, error) {
	return fuel.FuelEvent{ID: id, CarID: "1", Volume: 40, Price: 2000, Type: "AI-95", Date: "2026-03-22"}, nil
}
func (m *mockFuelRepo) UpdateFuelEvent(e fuel.FuelEvent) error { return nil }
func (m *mockFuelRepo) DeleteFuelEvent(id string) error        { return nil }
func (m *mockFuelRepo) ListFuelEvents() ([]fuel.FuelEvent, error) {
	return []fuel.FuelEvent{{ID: "1", CarID: "1", Volume: 40, Price: 2000, Type: "AI-95", Date: "2026-03-22"}}, nil
}

func TestAddFuelEventUsecase_Execute(t *testing.T) {
	uc := usecase.AddFuelEventUsecase{Repo: &mockFuelRepo{}}
	event := fuel.FuelEvent{ID: "1", CarID: "1", Volume: 40, Price: 2000, Type: "AI-95", Date: "2026-03-22"}
	err := uc.Execute(event)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestGetFuelEventUsecase_Execute(t *testing.T) {
	uc := usecase.GetFuelEventUsecase{Repo: &mockFuelRepo{}}
	_, err := uc.Execute("1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestUpdateFuelEventUsecase_Execute(t *testing.T) {
	uc := usecase.UpdateFuelEventUsecase{Repo: &mockFuelRepo{}}
	event := fuel.FuelEvent{ID: "1", CarID: "1", Volume: 50, Price: 2100, Type: "AI-95", Date: "2026-03-23"}
	err := uc.Execute(event)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestDeleteFuelEventUsecase_Execute(t *testing.T) {
	uc := usecase.DeleteFuelEventUsecase{Repo: &mockFuelRepo{}}
	err := uc.Execute("1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestListFuelEventsUsecase_Execute(t *testing.T) {
	uc := usecase.ListFuelEventsUsecase{Repo: &mockFuelRepo{}}
	events, err := uc.Execute()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(events) == 0 {
		t.Errorf("expected non-empty list")
	}
}

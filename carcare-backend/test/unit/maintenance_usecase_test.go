package unit

import (
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type mockMaintenanceRepo struct{}

func (m *mockMaintenanceRepo) AddMaintenanceEvent(e maintenance.MaintenanceEvent) error {
	return nil
}
func (m *mockMaintenanceRepo) GetMaintenanceEvent(id string) (maintenance.MaintenanceEvent, error) {
	return maintenance.MaintenanceEvent{}, nil
}
func (m *mockMaintenanceRepo) UpdateMaintenanceEvent(e maintenance.MaintenanceEvent) error {
	return nil
}
func (m *mockMaintenanceRepo) DeleteMaintenanceEvent(id string) error { return nil }
func (m *mockMaintenanceRepo) ListMaintenanceEvents(userID string) ([]maintenance.MaintenanceEvent, error) {
	return []maintenance.MaintenanceEvent{}, nil
}

func TestMaintenanceUsecase_Execute(t *testing.T) {
	uc := usecase.AddMaintenanceEventUsecase{Repo: &mockMaintenanceRepo{}}
	event := maintenance.MaintenanceEvent{ID: "1", CarID: "1", Type: "oil change", Date: "2026-03-22", Cost: 500}
	err := uc.Execute(event)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

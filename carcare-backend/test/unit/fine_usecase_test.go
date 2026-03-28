package unit

import (
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type mockFineRepo struct{}

func (m *mockFineRepo) AddFine(f fine.Fine) error                    { return nil }
func (m *mockFineRepo) GetFine(id string) (fine.Fine, error)         { return fine.Fine{}, nil }
func (m *mockFineRepo) UpdateFine(f fine.Fine) error                 { return nil }
func (m *mockFineRepo) DeleteFine(id string) error                   { return nil }
func (m *mockFineRepo) ListFines(userID string) ([]fine.Fine, error) { return []fine.Fine{}, nil }

func TestAddFineUsecase_Execute(t *testing.T) {
	uc := usecase.AddFineUsecase{Repo: &mockFineRepo{}}
	fineObj := fine.Fine{ID: "1", CarID: "1", Amount: 1000, Type: "speeding", Date: "2026-03-22"}
	err := uc.Execute(fineObj)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

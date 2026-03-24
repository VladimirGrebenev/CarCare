package unit

import (
	"context"
	"testing"

	repo "github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
)

func TestFineRepository_GetByID_NotImplemented(t *testing.T) {
	r := repo.NewFineRepository()
	_, err := r.GetByID(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

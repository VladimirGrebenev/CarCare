package unit

import (
	"context"
	"testing"

	repo "github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
)

func TestCarRepository_GetByID_NotImplemented(t *testing.T) {
	r := repo.NewCarRepository()
	_, err := r.GetByID(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

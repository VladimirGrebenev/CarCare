package unit

import (
	"context"
	"testing"

	repo "home/veneberg/CarCare/carcare-backend/internal/adapter/repository"
)

func TestUserRepository_GetByID_NotImplemented(t *testing.T) {
	r := repo.NewUserRepository()
	_, err := r.GetByID(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

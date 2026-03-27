package unit

import (
	"testing"

	repo "github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
)

func TestReportRepository_New(t *testing.T) {
	r := repo.NewReportRepository(nil)
	if r == nil {
		t.Fatal("expected non-nil repository")
	}
}

package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
)

func newFineRepoWithMock() (*repository.FineRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	cleanup := func() { db.Close() }
	return repository.NewFineRepository(db), mock, cleanup
}

func TestFineRepository_AddFine(t *testing.T) {
	repo, mock, cleanup := newFineRepoWithMock()
	defer cleanup()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fines (id, car_id, amount, type, date) VALUES ($1, $2, $3, $4, $5)")).
		WithArgs("1", "1", 500.0, "speeding", "2026-03-22").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.AddFine(fine.Fine{ID: "1", CarID: "1", Amount: 500, Type: "speeding", Date: "2026-03-22"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFineRepository_ListFines(t *testing.T) {
	repo, mock, cleanup := newFineRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "amount", "type", "date"}).
		AddRow("1", "1", 500.0, "speeding", "2026-03-22")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, car_id, amount, type, date FROM fines")).
		WillReturnRows(rows)
	fines, err := repo.ListFines()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(fines) == 0 {
		t.Errorf("expected non-empty list")
	}
}

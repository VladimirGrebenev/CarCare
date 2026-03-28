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
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fines (id, car_id, amount, type, date, status, description) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
		WithArgs("1", "1", 500.0, "speeding", "2026-03-22", "unpaid", "превышение скорости").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.AddFine(fine.Fine{ID: "1", CarID: "1", Amount: 500, Type: "speeding", Date: "2026-03-22", Status: "unpaid", Description: "превышение скорости"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFineRepository_ListFines(t *testing.T) {
	repo, mock, cleanup := newFineRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "amount", "type", "date", "status", "description"}).
		AddRow("1", "1", 500.0, "speeding", "2026-03-22", "unpaid", "превышение скорости")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT f.id, f.car_id, f.amount, f.type, f.date, f.status, f.description FROM fines f JOIN cars c ON c.id = f.car_id WHERE c.user_id = $1")).
		WithArgs("user-1").
		WillReturnRows(rows)
	fines, err := repo.ListFines("user-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(fines) == 0 {
		t.Errorf("expected non-empty list")
	}
}

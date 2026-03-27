package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
)

func newMaintenanceRepoWithMock() (*repository.MaintenanceRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	cleanup := func() { db.Close() }
	return repository.NewMaintenanceRepository(db), mock, cleanup
}

func TestMaintenanceRepository_AddMaintenanceEvent(t *testing.T) {
	repo, mock, cleanup := newMaintenanceRepoWithMock()
	defer cleanup()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO maintenance_events (id, car_id, type, date, cost) VALUES ($1, $2, $3, $4, $5)")).
		WithArgs("1", "1", "oil_change", "2026-03-22", 3000.0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.AddMaintenanceEvent(maintenance.MaintenanceEvent{ID: "1", CarID: "1", Type: "oil_change", Date: "2026-03-22", Cost: 3000})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMaintenanceRepository_ListMaintenanceEvents(t *testing.T) {
	repo, mock, cleanup := newMaintenanceRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "type", "date", "cost"}).
		AddRow("1", "1", "oil_change", "2026-03-22", 3000.0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, car_id, type, date, cost FROM maintenance_events")).
		WillReturnRows(rows)
	events, err := repo.ListMaintenanceEvents()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Errorf("expected non-empty list")
	}
}

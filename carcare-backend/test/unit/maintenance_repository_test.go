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
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO maintenance_events (id, car_id, type, date, cost, description) VALUES ($1, $2, $3, $4, $5, $6)")).
		WithArgs("1", "1", "oil_change", "2026-03-22", 3000.0, "плановая замена масла").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.AddMaintenanceEvent(maintenance.MaintenanceEvent{ID: "1", CarID: "1", Type: "oil_change", Date: "2026-03-22", Cost: 3000, Description: "плановая замена масла"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMaintenanceRepository_ListMaintenanceEvents(t *testing.T) {
	repo, mock, cleanup := newMaintenanceRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "type", "date", "cost", "description"}).
		AddRow("1", "1", "oil_change", "2026-03-22", 3000.0, "плановая замена масла")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT me.id, me.car_id, me.type, me.date, me.cost, me.description FROM maintenance_events me JOIN cars c ON c.id = me.car_id WHERE c.user_id = $1")).
		WithArgs("user-1").
		WillReturnRows(rows)
	events, err := repo.ListMaintenanceEvents("user-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Errorf("expected non-empty list")
	}
}

package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
)

func newFuelRepoWithMock() (*repository.FuelRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	cleanup := func() { db.Close() }
	return repository.NewFuelRepository(db), mock, cleanup
}

func TestFuelRepository_AddFuelEvent(t *testing.T) {
	repo, mock, cleanup := newFuelRepoWithMock()
	defer cleanup()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fuel_events (id, car_id, volume, price, type, date) VALUES ($1, $2, $3, $4, $5, $6)")).
		WithArgs("1", "1", 40.0, 2000.0, "AI-95", "2026-03-22").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.AddFuelEvent(fuel.FuelEvent{ID: "1", CarID: "1", Volume: 40, Price: 2000, Type: "AI-95", Date: "2026-03-22"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFuelRepository_GetFuelEvent(t *testing.T) {
	repo, mock, cleanup := newFuelRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "volume", "price", "type", "date"}).
		AddRow("1", "1", 40.0, 2000.0, "AI-95", "2026-03-22")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, car_id, volume, price, type, date FROM fuel_events WHERE id = $1")).
		WithArgs("1").WillReturnRows(rows)
	_, err := repo.GetFuelEvent("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFuelRepository_UpdateFuelEvent(t *testing.T) {
	repo, mock, cleanup := newFuelRepoWithMock()
	defer cleanup()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE fuel_events SET car_id=$1, volume=$2, price=$3, type=$4, date=$5 WHERE id=$6")).
		WithArgs("1", 50.0, 2100.0, "AI-95", "2026-03-23", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.UpdateFuelEvent(fuel.FuelEvent{ID: "1", CarID: "1", Volume: 50, Price: 2100, Type: "AI-95", Date: "2026-03-23"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFuelRepository_DeleteFuelEvent(t *testing.T) {
	repo, mock, cleanup := newFuelRepoWithMock()
	defer cleanup()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fuel_events WHERE id=$1")).
		WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.DeleteFuelEvent("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFuelRepository_ListFuelEvents(t *testing.T) {
	repo, mock, cleanup := newFuelRepoWithMock()
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "car_id", "volume", "price", "type", "date"}).
		AddRow("1", "1", 40.0, 2000.0, "AI-95", "2026-03-22")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT fe.id, fe.car_id, fe.volume, fe.price, fe.type, fe.date FROM fuel_events fe JOIN cars c ON c.id = fe.car_id WHERE c.user_id = $1")).
		WithArgs("user-1").
		WillReturnRows(rows)
	events, err := repo.ListFuelEvents("user-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Errorf("expected non-empty list")
	}
}

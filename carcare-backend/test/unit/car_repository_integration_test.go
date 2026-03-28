package unit

import (
	"database/sql"
	"os"
	"testing"

	repo "github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("CARCARE_TEST_DSN")
	if dsn == "" {
		t.Skip("CARCARE_TEST_DSN not set, skipping integration test")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	cleanup := func() {
		db.Exec("DELETE FROM cars")
	}
	cleanup()
	t.Cleanup(cleanup)
	return db
}

func TestCarRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewCarRepository(db)
	c := car.Car{ID: "uuid-1", UserID: "user-uuid-1", Brand: "Toyota", Model: "Corolla", Year: 2020, VIN: "VIN-1"}
	// Create
	if err := r.AddCar(c); err != nil {
		t.Fatalf("AddCar: %v", err)
	}
	// Read
	got, err := r.GetCar(c.ID)
	if err != nil || got.VIN != c.VIN {
		t.Fatalf("GetCar: %v, got: %+v", err, got)
	}
	// Update
	c.Brand = "Honda"
	if err := r.UpdateCar(c, c.UserID); err != nil {
		t.Fatalf("UpdateCar: %v", err)
	}
	// List
	cars, err := r.ListCars(c.UserID)
	if err != nil || len(cars) == 0 {
		t.Fatalf("ListCars: %v", err)
	}
	// Delete
	if err := r.DeleteCar(c.ID, c.UserID); err != nil {
		t.Fatalf("DeleteCar: %v", err)
	}
}

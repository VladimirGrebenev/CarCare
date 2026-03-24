package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/repository"
	"github.com/VladimirGrebenev/CarCare-backend/internal/adapter/rest"
	"github.com/VladimirGrebenev/CarCare-backend/internal/infrastructure/db"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func main() {
	fmt.Println("CarCare backend skeleton running...")

	// DI: инициализация репозиториев и usecase
	dsn := os.Getenv("CARCARE_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/carcare?sslmode=disable"
	}
	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		panic(err)
	}
	carRepo := repository.NewCarRepository(dbConn)
	fuelRepo := repository.NewFuelRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)
	uc := &usecase.UsecaseContainer{
		Car:         carRepo,
		Fuel:        fuelRepo,
		User:        userRepo,
		UserService: usecase.NewUserService(userRepo),
	}
	carHandler := rest.NewCarHandler(uc)
	userHandler := rest.NewUserHandler(uc)
	fuelHandler := rest.NewFuelHandler(uc)

	http.HandleFunc("/health", rest.HealthCheckHandler)
	http.Handle("/cars", carHandler)
	http.Handle("/cars/", carHandler) // для /cars/{id}
	http.Handle("/users", userHandler)
	http.Handle("/users/", userHandler)
	http.Handle("/fuel", fuelHandler)
	http.Handle("/fuel/", fuelHandler)
	http.HandleFunc("/maintenance", rest.MaintenanceHandler)
	http.HandleFunc("/fines", rest.FineHandler)
	http.HandleFunc("/reports", rest.ReportHandler)
	http.HandleFunc("/api/profile", rest.ProfileHandler)

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

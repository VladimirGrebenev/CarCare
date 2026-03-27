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
	maintenanceRepo := repository.NewMaintenanceRepository(dbConn)
	fineRepo := repository.NewFineRepository(dbConn)
	uc := usecase.NewUsecaseContainer(carRepo, fuelRepo, fineRepo, maintenanceRepo, userRepo)
	authUC := usecase.NewAuthUsecase(userRepo, nil, nil, nil, nil)
	carHandler := rest.NewCarHandler(uc)
	userHandler := rest.NewUserHandler(uc)
	fuelHandler := rest.NewFuelHandler(uc)
	maintenanceHandler := rest.NewMaintenanceHandler(uc)
	fineHandler := rest.NewFineHandler(uc)
	reportHandler := rest.NewReportHandler(uc)
	authHandler := rest.NewAuthHandler(authUC)

	// Публичные маршруты — без аутентификации
	http.HandleFunc("/health", rest.HealthCheckHandler)
	http.HandleFunc("/api/auth/login", authHandler.Login)
	http.HandleFunc("/api/auth/register", authHandler.Register)
	http.HandleFunc("/api/auth/oauth/", authHandler.OAuthProvider)

	// Защищённые маршруты — требуют валидный JWT
	http.Handle("/cars", rest.AuthMiddleware(carHandler))
	http.Handle("/cars/", rest.AuthMiddleware(carHandler))
	http.Handle("/api/cars", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", carHandler)))
	http.Handle("/api/cars/", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", carHandler)))
	http.Handle("/users", rest.AuthMiddleware(userHandler))
	http.Handle("/users/", rest.AuthMiddleware(userHandler))
	http.Handle("/fuel", rest.AuthMiddleware(fuelHandler))
	http.Handle("/fuel/", rest.AuthMiddleware(fuelHandler))
	http.Handle("/api/fuel", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", fuelHandler)))
	http.Handle("/api/fuel/", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", fuelHandler)))
	http.Handle("/maintenance", rest.AuthMiddleware(maintenanceHandler))
	http.Handle("/maintenance/", rest.AuthMiddleware(maintenanceHandler))
	http.Handle("/api/maintenance", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", maintenanceHandler)))
	http.Handle("/api/maintenance/", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", maintenanceHandler)))
	http.Handle("/fines", rest.AuthMiddleware(fineHandler))
	http.Handle("/fines/", rest.AuthMiddleware(fineHandler))
	http.Handle("/api/fines", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", fineHandler)))
	http.Handle("/api/fines/", rest.AuthMiddleware(rest.AliasPrefixHandler("/api", "", fineHandler)))
	http.Handle("/reports", rest.AuthMiddleware(reportHandler))
	http.Handle("/api/reports", rest.AuthMiddleware(reportHandler))
	http.Handle("/api/reports/export", rest.AuthMiddleware(reportHandler))
	http.Handle("/api/profile", rest.AuthMiddleware(http.HandlerFunc(rest.ProfileHandler)))

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

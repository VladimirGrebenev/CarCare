package usecase

// Dependency injection container for usecases

import (
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
)

type UsecaseContainer struct {
	Car         car.Repository
	Fuel        fuel.Repository
	Fine        fine.Repository
	Maintenance maintenance.Repository
	User        user.Repository
	UserService UserService
}

func NewUsecaseContainer(
	carRepo car.Repository,
	fuelRepo fuel.Repository,
	fineRepo fine.Repository,
	maintenanceRepo maintenance.Repository,
	userRepo user.Repository,
) *UsecaseContainer {
	return &UsecaseContainer{
		Car:         carRepo,
		Fuel:        fuelRepo,
		Fine:        fineRepo,
		Maintenance: maintenanceRepo,
		User:        userRepo,
		UserService: NewUserService(userRepo),
	}
}

package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
)

type AddCarUsecase struct {
	Repo car.Repository
}

// Валидация и доменные правила
func validateCar(c car.Car) error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(c.Brand) == "" {
		return errors.New("brand is required")
	}
	if strings.TrimSpace(c.Model) == "" {
		return errors.New("model is required")
	}
	if c.Year < 1886 || c.Year > 2100 {
		return fmt.Errorf("year %d is invalid", c.Year)
	}
	if strings.TrimSpace(c.VIN) == "" {
		return errors.New("VIN is required")
	}
	// TODO: добавить доменные правила (например, уникальность VIN)
	return nil
}

func (uc *AddCarUsecase) Execute(c car.Car) error {
	if err := validateCar(c); err != nil {
		return err
	}
	// TODO: проверить уникальность VIN через Repo
	return uc.Repo.AddCar(c)
}

type GetCarUsecase struct {
	Repo car.Repository
}

func (uc *GetCarUsecase) Execute(id string) (car.Car, error) {
	if strings.TrimSpace(id) == "" {
		return car.Car{}, errors.New("id is required")
	}
	return uc.Repo.GetCar(id)
}

type UpdateCarUsecase struct {
	Repo car.Repository
}

func (uc *UpdateCarUsecase) Execute(c car.Car, userID string) error {
	if err := validateCar(c); err != nil {
		return err
	}
	return uc.Repo.UpdateCar(c, userID)
}

type DeleteCarUsecase struct {
	Repo car.Repository
}

func (uc *DeleteCarUsecase) Execute(id string, userID string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	return uc.Repo.DeleteCar(id, userID)
}

type ListCarsUsecase struct {
	Repo car.Repository
}

func (uc *ListCarsUsecase) Execute(userID string) ([]car.Car, error) {
	return uc.Repo.ListCars(userID)
}

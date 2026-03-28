package unit

import (
	"testing"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

func TestAddCarUsecase_Execute(t *testing.T) {
	repo := NewMockCarRepo()
	uc := usecase.AddCarUsecase{Repo: repo}
	carObj := car.Car{ID: "1", Brand: "TestBrand", Model: "TestModel", Year: 2020, VIN: "VIN123"}
	err := uc.Execute(carObj)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	// Проверка, что машина добавлена
	got, err := repo.GetCar("1")
	if err != nil || got.ID != "1" {
		t.Errorf("car not added: %v", err)
	}
}

func TestGetCarUsecase_Execute(t *testing.T) {
	repo := NewMockCarRepo()
	repo.AddCar(car.Car{ID: "2", Brand: "B", Model: "M", Year: 2021, VIN: "VIN2"})
	uc := usecase.GetCarUsecase{Repo: repo}
	c, err := uc.Execute("2")
	if err != nil || c.ID != "2" {
		t.Errorf("expected car, got %v, err %v", c, err)
	}
}

func TestUpdateCarUsecase_Execute(t *testing.T) {
	repo := NewMockCarRepo()
	repo.AddCar(car.Car{ID: "3", Brand: "B", Model: "M", Year: 2021, VIN: "VIN3"})
	uc := usecase.UpdateCarUsecase{Repo: repo}
	updated := car.Car{ID: "3", Brand: "B2", Model: "M2", Year: 2022, VIN: "VIN3"}
	err := uc.Execute(updated, "user-1")
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	got, _ := repo.GetCar("3")
	if got.Brand != "B2" || got.Model != "M2" || got.Year != 2022 {
		t.Errorf("update not applied: %+v", got)
	}
}

func TestDeleteCarUsecase_Execute(t *testing.T) {
	repo := NewMockCarRepo()
	repo.AddCar(car.Car{ID: "4", Brand: "B", Model: "M", Year: 2021, VIN: "VIN4"})
	uc := usecase.DeleteCarUsecase{Repo: repo}
	err := uc.Execute("4", "user-1")
	if err != nil {
		t.Errorf("delete failed: %v", err)
	}
	_, err = repo.GetCar("4")
	if err == nil {
		t.Errorf("car not deleted")
	}
}

func TestListCarsUsecase_Execute(t *testing.T) {
	repo := NewMockCarRepo()
	repo.AddCar(car.Car{ID: "5", Brand: "B", Model: "M", Year: 2021, VIN: "VIN5"})
	repo.AddCar(car.Car{ID: "6", Brand: "B2", Model: "M2", Year: 2022, VIN: "VIN6"})
	uc := usecase.ListCarsUsecase{Repo: repo}
	cars, err := uc.Execute("user-1")
	if err != nil || len(cars) != 2 {
		t.Errorf("expected 2 cars, got %d, err %v", len(cars), err)
	}
}

package usecase

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"


type AddFuelEventUsecase struct {
	Repo fuel.Repository
}

func (uc *AddFuelEventUsecase) Execute(e fuel.FuelEvent) error {
	// TODO: валидация
	return uc.Repo.AddFuelEvent(e)
}

type GetFuelEventUsecase struct {
	Repo fuel.Repository
}

func (uc *GetFuelEventUsecase) Execute(id string) (fuel.FuelEvent, error) {
	return uc.Repo.GetFuelEvent(id)
}

type UpdateFuelEventUsecase struct {
	Repo fuel.Repository
}

func (uc *UpdateFuelEventUsecase) Execute(e fuel.FuelEvent) error {
	// TODO: валидация
	return uc.Repo.UpdateFuelEvent(e)
}

type DeleteFuelEventUsecase struct {
	Repo fuel.Repository
}

func (uc *DeleteFuelEventUsecase) Execute(id string) error {
	return uc.Repo.DeleteFuelEvent(id)
}

type ListFuelEventsUsecase struct {
	Repo fuel.Repository
}

func (uc *ListFuelEventsUsecase) Execute() ([]fuel.FuelEvent, error) {
	return uc.Repo.ListFuelEvents()
}

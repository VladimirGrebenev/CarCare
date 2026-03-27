package usecase

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"

type AddMaintenanceEventUsecase struct {
	Repo maintenance.Repository
}

func (uc *AddMaintenanceEventUsecase) Execute(e maintenance.MaintenanceEvent) error {
	return uc.Repo.AddMaintenanceEvent(e)
}

type GetMaintenanceEventUsecase struct {
	Repo maintenance.Repository
}

func (uc *GetMaintenanceEventUsecase) Execute(id string) (maintenance.MaintenanceEvent, error) {
	return uc.Repo.GetMaintenanceEvent(id)
}

type UpdateMaintenanceEventUsecase struct {
	Repo maintenance.Repository
}

func (uc *UpdateMaintenanceEventUsecase) Execute(e maintenance.MaintenanceEvent) error {
	return uc.Repo.UpdateMaintenanceEvent(e)
}

type DeleteMaintenanceEventUsecase struct {
	Repo maintenance.Repository
}

func (uc *DeleteMaintenanceEventUsecase) Execute(id string) error {
	return uc.Repo.DeleteMaintenanceEvent(id)
}

type ListMaintenanceEventsUsecase struct {
	Repo maintenance.Repository
}

func (uc *ListMaintenanceEventsUsecase) Execute() ([]maintenance.MaintenanceEvent, error) {
	return uc.Repo.ListMaintenanceEvents()
}

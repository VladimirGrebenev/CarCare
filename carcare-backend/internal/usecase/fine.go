package usecase

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"

type AddFineUsecase struct {
	Repo fine.Repository
}

func (uc *AddFineUsecase) Execute(f fine.Fine) error {
	return uc.Repo.AddFine(f)
}

type GetFineUsecase struct {
	Repo fine.Repository
}

func (uc *GetFineUsecase) Execute(id string) (fine.Fine, error) {
	return uc.Repo.GetFine(id)
}

type UpdateFineUsecase struct {
	Repo fine.Repository
}

func (uc *UpdateFineUsecase) Execute(f fine.Fine) error {
	return uc.Repo.UpdateFine(f)
}

type DeleteFineUsecase struct {
	Repo fine.Repository
}

func (uc *DeleteFineUsecase) Execute(id string) error {
	return uc.Repo.DeleteFine(id)
}

type ListFinesUsecase struct {
	Repo fine.Repository
}

func (uc *ListFinesUsecase) Execute(userID string) ([]fine.Fine, error) {
	return uc.Repo.ListFines(userID)
}

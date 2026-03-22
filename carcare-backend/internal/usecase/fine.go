package usecase

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"

type AddFineUsecase struct {
	Repo fine.Repository
}

func (uc *AddFineUsecase) Execute(f fine.Fine) error {
	// TODO: implement
	return nil
}

package usecase

import "github.com/VladimirGrebenev/CarCare-backend/internal/domain/report"

type GetReportUsecase struct {
	Repo report.Repository
}

func (uc *GetReportUsecase) Execute(id string) (*report.Report, error) {
	return uc.Repo.GetByID(id)
}
package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
)

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

// GosuslugiImportBill — данные одного штрафа из внешнего источника
type GosuslugiImportBill struct {
	BillNumber string
	BillDate   int64   // Unix milliseconds
	BillName   string
	Amount     float64
	IsPaid     bool
}

// GosuslugiFetcher — интерфейс получения штрафов из внешнего источника
type GosuslugiFetcher interface {
	FetchBySts(sts string) ([]GosuslugiImportBill, error)
}

// ImportResult — результат импорта штрафов
type ImportResult struct {
	Added   int `json:"added"`
	Skipped int `json:"skipped"`
}

// ImportFinesByStsUsecase — получает штрафы по СТС и сохраняет только новые
type ImportFinesByStsUsecase struct {
	Repo fine.Repository
}

func (uc *ImportFinesByStsUsecase) Execute(carID string, fetcher GosuslugiFetcher, sts string) (ImportResult, error) {
	bills, err := fetcher.FetchBySts(sts)
	if err != nil {
		return ImportResult{}, err
	}

	var result ImportResult
	for _, bill := range bills {
		if bill.BillNumber != "" {
			exists, err := uc.Repo.CheckFineExistsByBillNumber(carID, bill.BillNumber)
			if err != nil {
				return result, err
			}
			if exists {
				result.Skipped++
				continue
			}
		}

		billDate := time.UnixMilli(bill.BillDate).Format("2006-01-02")
		status := "unpaid"
		if bill.IsPaid {
			status = "paid"
		}

		f := fine.Fine{
			ID:          uuid.New().String(),
			CarID:       carID,
			Amount:      bill.Amount,
			Type:        "Госуслуги",
			Date:        billDate,
			Status:      status,
			Description: bill.BillName,
			BillNumber:  bill.BillNumber,
		}

		if err := uc.Repo.AddFine(f); err != nil {
			return result, err
		}
		result.Added++
	}

	return result, nil
}

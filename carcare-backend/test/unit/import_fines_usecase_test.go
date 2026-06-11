package unit

import (
	"testing"
	"time"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type stubFetcher struct {
	bills []usecase.GosuslugiImportBill
	err   error
}

func (s *stubFetcher) FetchBySts(_ string) ([]usecase.GosuslugiImportBill, error) {
	return s.bills, s.err
}

// mockFineRepoReturnsExisting — CheckFineExistsByBillNumber возвращает true для заданного billNumber
type mockFineRepoReturnsExisting struct {
	existingBillNumber string
}

func (m *mockFineRepoReturnsExisting) AddFine(_ fine.Fine) error                   { return nil }
func (m *mockFineRepoReturnsExisting) GetFine(_ string) (fine.Fine, error)         { return fine.Fine{}, nil }
func (m *mockFineRepoReturnsExisting) UpdateFine(_ fine.Fine) error                { return nil }
func (m *mockFineRepoReturnsExisting) DeleteFine(_ string) error                   { return nil }
func (m *mockFineRepoReturnsExisting) ListFines(_ string) ([]fine.Fine, error)     { return []fine.Fine{}, nil }
func (m *mockFineRepoReturnsExisting) CheckFineExistsByBillNumber(_ string, billNumber string) (bool, error) {
	return billNumber == m.existingBillNumber, nil
}

func billsFixture() []usecase.GosuslugiImportBill {
	return []usecase.GosuslugiImportBill{
		{BillNumber: "BILL001", BillDate: time.Now().UnixMilli(), Amount: 500, BillName: "Штраф 1", IsPaid: false},
		{BillNumber: "BILL002", BillDate: time.Now().UnixMilli(), Amount: 300, BillName: "Штраф 2", IsPaid: true},
	}
}

func TestImportFinesBySts_AddsNewFines(t *testing.T) {
	uc := usecase.ImportFinesByStsUsecase{Repo: &mockFineRepo{}}
	fetcher := &stubFetcher{bills: billsFixture()}

	result, err := uc.Execute("car-1", fetcher, "99АА123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Added != 2 {
		t.Errorf("expected 2 added, got %d", result.Added)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", result.Skipped)
	}
}

func TestImportFinesBySts_SkipsDuplicates(t *testing.T) {
	repo := &mockFineRepoReturnsExisting{existingBillNumber: "BILL001"}
	uc := usecase.ImportFinesByStsUsecase{Repo: repo}
	fetcher := &stubFetcher{bills: billsFixture()}

	result, err := uc.Execute("car-1", fetcher, "99АА123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

func TestImportFinesBySts_EmptyBills_ReturnsZero(t *testing.T) {
	uc := usecase.ImportFinesByStsUsecase{Repo: &mockFineRepo{}}
	fetcher := &stubFetcher{bills: []usecase.GosuslugiImportBill{}}

	result, err := uc.Execute("car-1", fetcher, "99АА123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Added != 0 || result.Skipped != 0 {
		t.Errorf("expected 0/0, got %d/%d", result.Added, result.Skipped)
	}
}

func TestImportFinesBySts_EmptyBillNumber_IsSkipped(t *testing.T) {
	uc := usecase.ImportFinesByStsUsecase{Repo: &mockFineRepo{}}
	fetcher := &stubFetcher{bills: []usecase.GosuslugiImportBill{
		{BillNumber: "", BillDate: time.Now().UnixMilli(), Amount: 100, BillName: "Без номера", IsPaid: false},
	}}

	result, err := uc.Execute("car-1", fetcher, "99АА123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Added != 0 {
		t.Errorf("expected 0 added, got %d", result.Added)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

func TestImportFinesBySts_ZeroBillDate_IsSkipped(t *testing.T) {
	uc := usecase.ImportFinesByStsUsecase{Repo: &mockFineRepo{}}
	fetcher := &stubFetcher{bills: []usecase.GosuslugiImportBill{
		{BillNumber: "BILL001", BillDate: 0, Amount: 500, BillName: "Нет даты", IsPaid: false},
	}}

	result, err := uc.Execute("car-1", fetcher, "99АА123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Added != 0 {
		t.Errorf("expected 0 added, got %d", result.Added)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

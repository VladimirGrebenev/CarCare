package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// GosuslugiFineRequest — запрос к API Госуслуг
type GosuslugiFineRequest struct {
	Documents []GosuslugiDocument `json:"documents"`
}

type GosuslugiDocument struct {
	DocType   string `json:"docType"`
	DocNumber string `json:"docNumber"`
}

// GosuslugiFineResponse — ответ от API Госуслуг
type GosuslugiFineResponse struct {
	Error  GosuslugiError   `json:"error"`
	Errors []GosuslugiError `json:"errors"`
	Bills  []GosuslugiBill  `json:"bills"`
}

type GosuslugiError struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type GosuslugiBill struct {
	BillID     int64   `json:"billId"`
	BillNumber string  `json:"billNumber"`
	BillName   string  `json:"billName"`
	BillDate   int64   `json:"billDate"`
	IsPaid     bool    `json:"isPaid"`
	Amount     float64 `json:"amount"`
}

// CheckFinesByStsRequest — запрос от фронтенда
type CheckFinesByStsRequest struct {
	Sts string `json:"sts"`
}

// CheckFinesByStsResponse — ответ фронтенду
type CheckFinesByStsResponse struct {
	Fines []GosuslugiBill `json:"fines"`
	Error string          `json:"error,omitempty"`
}

// FinesByStsHandler — обрабатывает запросы на получение штрафов по СТС
type FinesByStsHandler struct {
	fineRepo fine.Repository
}

func NewFinesByStsHandler(uc *usecase.UsecaseContainer) *FinesByStsHandler {
	return &FinesByStsHandler{
		fineRepo: uc.Fine,
	}
}

func (h *FinesByStsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	var req CheckFinesByStsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.Sts == "" {
		http.Error(w, `{"error":"STS is required"}`, http.StatusBadRequest)
		return
	}

	fines, err := h.checkFines(req.Sts)
	if err != nil {
		json.NewEncoder(w).Encode(CheckFinesByStsResponse{
			Fines: []GosuslugiBill{},
			Error: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(CheckFinesByStsResponse{
		Fines: fines,
	})
}

func (h *FinesByStsHandler) checkFines(sts string) ([]GosuslugiBill, error) {
	url := "https://www.gosuslugi.ru/api/pay/public/v1/paygate/bill/create?serviceCategory=FINE&interfaceTypeCode=BETA_NOAUTH"

	payload := GosuslugiFineRequest{
		Documents: []GosuslugiDocument{
			{
				DocType:   "CAR_REG_CERTIFICATE",
				DocNumber: sts,
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Origin", "https://www.gosuslugi.ru")
	httpReq.Header.Set("Referer", "https://www.gosuslugi.ru/")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("сервис Госуслуг временно недоступен: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	var gosResp GosuslugiFineResponse
	if err := json.Unmarshal(respBody, &gosResp); err != nil {
		return nil, fmt.Errorf("ошибка при обработке ответа от Госуслуг: %v", err)
	}

	// Проверяем ошибки
	if gosResp.Error.ErrorCode != 0 {
		if gosResp.Error.ErrorCode == 1 {
			return []GosuslugiBill{}, nil // нет штрафов — не ошибка
		}
		return nil, fmt.Errorf("ошибка Госуслуг: %s", gosResp.Error.ErrorMessage)
	}

	if len(gosResp.Errors) > 0 && gosResp.Errors[0].ErrorCode != 0 {
		if gosResp.Errors[0].ErrorCode == 1 {
			return []GosuslugiBill{}, nil // нет штрафов
		}
		return nil, fmt.Errorf("ошибка Госуслуг: %s", gosResp.Errors[0].ErrorMessage)
	}

	return gosResp.Bills, nil
}

// CreateFineFromBill создаёт штраф из данных Госуслуг, проверяя дубликаты
func (h *FinesByStsHandler) CreateFineFromBill(carID string, bill GosuslugiBill) (bool, error) {
	// Проверяем, есть ли уже такой штраф
	if bill.BillNumber != "" {
		exists, err := h.fineRepo.CheckFineExistsByBillNumber(carID, bill.BillNumber)
		if err != nil {
			return false, err
		}
		if exists {
			return false, nil // уже есть — пропускаем
		}
	}

	billDate := time.UnixMilli(bill.BillDate).Format("2006-01-02")

	f := fine.Fine{
		ID:          fmt.Sprintf("%d", bill.BillID),
		CarID:       carID,
		Amount:      bill.Amount,
		Type:        "Госуслуги",
		Date:        billDate,
		Status:      "unpaid",
		Description: bill.BillName,
		BillNumber:  bill.BillNumber,
	}

	if bill.IsPaid {
		f.Status = "paid"
	}

	if err := h.fineRepo.AddFine(f); err != nil {
		return false, err
	}

	return true, nil
}
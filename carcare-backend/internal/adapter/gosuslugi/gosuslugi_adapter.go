package gosuslugi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// GosuslugiAdapter defines integration with the Gosuslugi external service.
type GosuslugiAdapter struct {
	httpClient *http.Client
}

func NewGosuslugiAdapter() *GosuslugiAdapter {
	return &GosuslugiAdapter{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// gosuslugiRequest — запрос к API Госуслуг
type gosuslugiRequest struct {
	Documents []gosuslugiDocument `json:"documents"`
}

type gosuslugiDocument struct {
	DocType   string `json:"docType"`
	DocNumber string `json:"docNumber"`
}

// gosuslugiResponse — ответ от API Госуслуг
type gosuslugiResponse struct {
	Error  struct {
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	} `json:"error"`
	Errors []struct {
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	} `json:"errors"`
	Bills []struct {
		BillID     int64   `json:"billId"`
		BillNumber string  `json:"billNumber"`
		BillName   string  `json:"billName"`
		BillDate   int64   `json:"billDate"`
		IsPaid     bool    `json:"isPaid"`
		Amount     float64 `json:"amount"`
	} `json:"bills"`
}

// FetchBySts реализует интерфейс usecase.GosuslugiFetcher.
// Делает HTTP-запрос к Госуслугам и возвращает список штрафов по СТС.
func (a *GosuslugiAdapter) FetchBySts(sts string) ([]usecase.GosuslugiImportBill, error) {
	const apiURL = "https://www.gosuslugi.ru/api/pay/public/v1/paygate/bill/create?serviceCategory=FINE&interfaceTypeCode=BETA_NOAUTH"

	payload := gosuslugiRequest{
		Documents: []gosuslugiDocument{
			{
				DocType:   "CAR_REG_CERTIFICATE",
				DocNumber: sts,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка при формировании запроса: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Origin", "https://www.gosuslugi.ru")
	httpReq.Header.Set("Referer", "https://www.gosuslugi.ru/")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("сервис Госуслуг временно недоступен: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервис Госуслуг вернул статус %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var gosResp gosuslugiResponse
	if err := json.Unmarshal(respBody, &gosResp); err != nil {
		return nil, fmt.Errorf("ошибка при обработке ответа от Госуслуг: %v", err)
	}

	// ErrorCode==1 означает "нет штрафов" — не ошибка
	if gosResp.Error.ErrorCode != 0 {
		if gosResp.Error.ErrorCode == 1 {
			return []usecase.GosuslugiImportBill{}, nil
		}
		return nil, fmt.Errorf("ошибка Госуслуг: %s", gosResp.Error.ErrorMessage)
	}

	if len(gosResp.Errors) > 0 && gosResp.Errors[0].ErrorCode != 0 {
		if gosResp.Errors[0].ErrorCode == 1 {
			return []usecase.GosuslugiImportBill{}, nil
		}
		return nil, fmt.Errorf("ошибка Госуслуг: %s", gosResp.Errors[0].ErrorMessage)
	}

	bills := make([]usecase.GosuslugiImportBill, 0, len(gosResp.Bills))
	for _, b := range gosResp.Bills {
		bills = append(bills, usecase.GosuslugiImportBill{
			BillNumber: b.BillNumber,
			BillDate:   b.BillDate,
			BillName:   b.BillName,
			Amount:     b.Amount,
			IsPaid:     b.IsPaid,
		})
	}

	return bills, nil
}

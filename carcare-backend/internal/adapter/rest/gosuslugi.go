package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

type ImportFinesByStsRequest struct {
	CarID string `json:"car_id"`
	Sts   string `json:"sts"`
}

type ImportFinesByStsResponse struct {
	Added   int    `json:"added"`
	Skipped int    `json:"skipped"`
	Error   string `json:"error,omitempty"`
}

type ImportFinesByStsHandler struct {
	uc      *usecase.ImportFinesByStsUsecase
	fetcher usecase.GosuslugiFetcher
}

func NewImportFinesByStsHandler(uc *usecase.UsecaseContainer, fetcher usecase.GosuslugiFetcher) *ImportFinesByStsHandler {
	return &ImportFinesByStsHandler{
		uc:      &usecase.ImportFinesByStsUsecase{Repo: uc.Fine},
		fetcher: fetcher,
	}
}

func (h *ImportFinesByStsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var req ImportFinesByStsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.CarID == "" || req.Sts == "" {
		http.Error(w, `{"error":"car_id and sts are required"}`, http.StatusBadRequest)
		return
	}

	result, err := h.uc.Execute(req.CarID, h.fetcher, req.Sts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ImportFinesByStsResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(ImportFinesByStsResponse{
		Added:   result.Added,
		Skipped: result.Skipped,
	})
}

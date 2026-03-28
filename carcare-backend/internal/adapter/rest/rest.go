package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fine"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/fuel"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/maintenance"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/user"
	"github.com/VladimirGrebenev/CarCare-backend/internal/usecase"
)

// CarHandler handles car-related requests
// CarHandler реализует CRUD для Car
type CarHandler struct {
	Add    *usecase.AddCarUsecase
	Get    *usecase.GetCarUsecase
	Update *usecase.UpdateCarUsecase
	Delete *usecase.DeleteCarUsecase
	List   *usecase.ListCarsUsecase
}

func NewCarHandler(uc *usecase.UsecaseContainer) *CarHandler {
	return &CarHandler{
		Add:    &usecase.AddCarUsecase{Repo: uc.Car},
		Get:    &usecase.GetCarUsecase{Repo: uc.Car},
		Update: &usecase.UpdateCarUsecase{Repo: uc.Car},
		Delete: &usecase.DeleteCarUsecase{Repo: uc.Car},
		List:   &usecase.ListCarsUsecase{Repo: uc.Car},
	}
}

func (h *CarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/cars/")
		if id == "" || id == "/cars" || r.URL.Path == "/cars" {
			h.handleList(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		h.handleAdd(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/cars/")
		h.handleDelete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CarHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	var c car.Car
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	c.UserID = userID
	if err := h.Add.Execute(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CarHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	c, err := h.Get.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func (h *CarHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	var c car.Car
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.Update.Execute(c, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *CarHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	userID := getUserIDFromContext(r)
	if err := h.Delete.Execute(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CarHandler) handleList(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	cars, err := h.List.Execute(userID)
	if err != nil || cars == nil {
		cars = []car.Car{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars)
}

// UserHandler handles user-related requests (CRUD)
type UserHandler struct {
	Service usecase.UserService
}

func NewUserHandler(uc *usecase.UsecaseContainer) *UserHandler {
	return &UserHandler{Service: uc.UserService}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" || id == "/users" || r.URL.Path == "/users" {
			h.handleList(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		h.handleDelete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var u usecaseUserDTO
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user := u.ToDomain()
	err = h.Service.Create(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	u, err := h.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var u usecaseUserDTO
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	user := u.ToDomain()
	err = h.Service.Update(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	err := h.Service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleList(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.List(r.Context())
	if err != nil || users == nil {
		users = []*user.User{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// DTO for user JSON
type usecaseUserDTO struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (dto usecaseUserDTO) ToDomain() user.User {
	return user.User{
		ID:    dto.ID,
		Email: user.Email(dto.Email),
		Name:  dto.Name,
		Role:  user.Role(dto.Role),
	}
}

// FuelHandler реализует CRUD для FuelEvent
type FuelHandler struct {
	Add    *usecase.AddFuelEventUsecase
	Get    *usecase.GetFuelEventUsecase
	Update *usecase.UpdateFuelEventUsecase
	Delete *usecase.DeleteFuelEventUsecase
	List   *usecase.ListFuelEventsUsecase
}

func NewFuelHandler(uc *usecase.UsecaseContainer) *FuelHandler {
	return &FuelHandler{
		Add:    &usecase.AddFuelEventUsecase{Repo: uc.Fuel},
		Get:    &usecase.GetFuelEventUsecase{Repo: uc.Fuel},
		Update: &usecase.UpdateFuelEventUsecase{Repo: uc.Fuel},
		Delete: &usecase.DeleteFuelEventUsecase{Repo: uc.Fuel},
		List:   &usecase.ListFuelEventsUsecase{Repo: uc.Fuel},
	}
}

func (h *FuelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/fuel/")
		if id == "" || id == "/fuel" || r.URL.Path == "/fuel" {
			h.handleList(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		h.handleAdd(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/fuel/")
		h.handleDelete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FuelHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var e fuel.FuelEvent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &e); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if err := h.Add.Execute(e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *FuelHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	e, err := h.Get.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(e)
}

func (h *FuelHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	idFromURL := strings.TrimPrefix(r.URL.Path, "/fuel/")
	if idFromURL == "" {
		http.Error(w, "missing id in URL", http.StatusBadRequest)
		return
	}
	var e fuel.FuelEvent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &e); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	e.ID = idFromURL
	if err := h.Update.Execute(e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
}

func (h *FuelHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Delete.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FuelHandler) handleList(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	events, err := h.List.Execute(userID)
	if err != nil || events == nil {
		events = []fuel.FuelEvent{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// MaintenanceHandler реализует CRUD для MaintenanceEvent
type MaintenanceHandler struct {
	Add    *usecase.AddMaintenanceEventUsecase
	Get    *usecase.GetMaintenanceEventUsecase
	Update *usecase.UpdateMaintenanceEventUsecase
	Delete *usecase.DeleteMaintenanceEventUsecase
	List   *usecase.ListMaintenanceEventsUsecase
}

func NewMaintenanceHandler(uc *usecase.UsecaseContainer) *MaintenanceHandler {
	return &MaintenanceHandler{
		Add:    &usecase.AddMaintenanceEventUsecase{Repo: uc.Maintenance},
		Get:    &usecase.GetMaintenanceEventUsecase{Repo: uc.Maintenance},
		Update: &usecase.UpdateMaintenanceEventUsecase{Repo: uc.Maintenance},
		Delete: &usecase.DeleteMaintenanceEventUsecase{Repo: uc.Maintenance},
		List:   &usecase.ListMaintenanceEventsUsecase{Repo: uc.Maintenance},
	}
}

func (h *MaintenanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/maintenance/")
		if id == "" || id == "/maintenance" || r.URL.Path == "/maintenance" {
			h.handleList(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		h.handleAdd(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/maintenance/")
		h.handleDelete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MaintenanceHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var e maintenance.MaintenanceEvent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &e); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if err := h.Add.Execute(e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *MaintenanceHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	e, err := h.Get.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(e)
}

func (h *MaintenanceHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	idFromURL := strings.TrimPrefix(r.URL.Path, "/maintenance/")
	if idFromURL == "" {
		http.Error(w, "missing id in URL", http.StatusBadRequest)
		return
	}
	var e maintenance.MaintenanceEvent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &e); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	e.ID = idFromURL
	if err := h.Update.Execute(e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
}

func (h *MaintenanceHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Delete.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MaintenanceHandler) handleList(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	events, err := h.List.Execute(userID)
	if err != nil || events == nil {
		events = []maintenance.MaintenanceEvent{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// FineHandler реализует CRUD для Fine
type FineHandler struct {
	Add    *usecase.AddFineUsecase
	Get    *usecase.GetFineUsecase
	Update *usecase.UpdateFineUsecase
	Delete *usecase.DeleteFineUsecase
	List   *usecase.ListFinesUsecase
}

func NewFineHandler(uc *usecase.UsecaseContainer) *FineHandler {
	return &FineHandler{
		Add:    &usecase.AddFineUsecase{Repo: uc.Fine},
		Get:    &usecase.GetFineUsecase{Repo: uc.Fine},
		Update: &usecase.UpdateFineUsecase{Repo: uc.Fine},
		Delete: &usecase.DeleteFineUsecase{Repo: uc.Fine},
		List:   &usecase.ListFinesUsecase{Repo: uc.Fine},
	}
}

func (h *FineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimPrefix(r.URL.Path, "/fines/")
		if id == "" || id == "/fines" || r.URL.Path == "/fines" {
			h.handleList(w, r)
		} else {
			h.handleGet(w, r, id)
		}
	case http.MethodPost:
		h.handleAdd(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/fines/")
		h.handleDelete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FineHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var f fine.Fine
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &f); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	if err := h.Add.Execute(f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

func (h *FineHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	f, err := h.Get.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(f)
}

func (h *FineHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	idFromURL := strings.TrimPrefix(r.URL.Path, "/fines/")
	if idFromURL == "" {
		http.Error(w, "missing id in URL", http.StatusBadRequest)
		return
	}
	var f fine.Fine
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &f); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	f.ID = idFromURL
	if err := h.Update.Execute(f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(f)
}

func (h *FineHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Delete.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FineHandler) handleList(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	fines, err := h.List.Execute(userID)
	if err != nil || fines == nil {
		fines = []fine.Fine{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fines)
}

// ReportSummary holds aggregated expense statistics
type ReportSummary struct {
	TotalFuelCost        float64 `json:"total_fuel_cost"`
	TotalMaintenanceCost float64 `json:"total_maintenance_cost"`
	TotalFinesAmount     float64 `json:"total_fines_amount"`
	FuelCount            int     `json:"fuel_count"`
	MaintenanceCount     int     `json:"maintenance_count"`
	FinesCount           int     `json:"fines_count"`
}

// ReportHandler handles report summary requests
type ReportHandler struct {
	UC *usecase.UsecaseContainer
}

func NewReportHandler(uc *usecase.UsecaseContainer) *ReportHandler {
	return &ReportHandler{UC: uc}
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := getUserIDFromContext(r)
	var summary ReportSummary
	if fuelEvents, err := (&usecase.ListFuelEventsUsecase{Repo: h.UC.Fuel}).Execute(userID); err == nil {
		summary.FuelCount = len(fuelEvents)
		for _, e := range fuelEvents {
			summary.TotalFuelCost += e.Volume * e.Price
		}
	}
	if mEvents, err := (&usecase.ListMaintenanceEventsUsecase{Repo: h.UC.Maintenance}).Execute(userID); err == nil {
		summary.MaintenanceCount = len(mEvents)
		for _, e := range mEvents {
			summary.TotalMaintenanceCost += e.Cost
		}
	}
	if fines, err := (&usecase.ListFinesUsecase{Repo: h.UC.Fine}).Execute(userID); err == nil {
		summary.FinesCount = len(fines)
		for _, f := range fines {
			summary.TotalFinesAmount += f.Amount
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

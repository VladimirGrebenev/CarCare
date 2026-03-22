package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/VladimirGrebenev/CarCare-backend/internal/domain/car"
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
	if err := h.Update.Execute(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *CarHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Delete.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CarHandler) handleList(w http.ResponseWriter, r *http.Request) {
	cars, err := h.List.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
       events, err := h.List.Execute()
       if err != nil {
	       http.Error(w, err.Error(), http.StatusInternalServerError)
	       return
       }
       json.NewEncoder(w).Encode(events)
}

// MaintenanceHandler handles maintenance-related requests
func MaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("maintenance endpoint not implemented"))
}

// FineHandler handles fine-related requests
func FineHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("fine endpoint not implemented"))
}

// ReportHandler handles report-related requests
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("report endpoint not implemented"))
}

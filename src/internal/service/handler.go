package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы для услуг.
type Handler struct {
	service *ServiceManager
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(service *ServiceManager) *Handler {
	return &Handler{service: service}
}

// GetServices обрабатывает GET /services.
func (h *Handler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.service.GetAllServices()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Преобразуем []Service в []ServiceResponse, если нужно скрыть поля
	resp := make([]ServiceResponse, len(services))
	for i, s := range services {
		resp[i] = ServiceResponse{ID: s.ID, Name: s.Name}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	service, err := h.service.CreateService(req.Name)
	if err != nil {
		// Если ошибка валидации – 400, иначе 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ServiceResponse{
		ID:   service.ID,
		Name: service.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// При успешном удалении возвращает статус 204 No Content.
// Если услуга не найдена, возвращает 404 Not Found.
func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL (предполагается использование chi)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteService(id)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		// Внутренняя ошибка сервера
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package admin

import (
	"encoding/json"
	"net/http"

	"src/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы для авторизации пользователей
type Handler struct {
	admin *AdminManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(service *AdminManager) *Handler {
	return &Handler{admin: service}
}

// CreatePartnerRequest обрабатывает POST /partner/request, создает заявку для организации
func (h *Handler) CreatePartnerRequest(w http.ResponseWriter, r *http.Request) {
	// Получение claims
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Извлечение email из claims
	userEmail, ok := claims["email"].(string)
	if !ok || userEmail == "" {
		http.Error(w, "Invalid token: email not found", http.StatusUnauthorized)
		return
	}

	var req PartnerRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.CreatePartnerRequest(userEmail, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Partner request created successfully",
	})
}

// TakeRequestToWork обрабатывает POST /admin/partner-requests/take, меняет статус заявки
func (h *Handler) TakeRequestToWork(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.TakeRequestToWork(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request taken to work",
		"id":      req.ID.String(),
		"status":  "pending",
	})
}

// GetAllRequests обрабатывает GET /admin/partner-requests/, получает все заявки
func (h *Handler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.admin.GetAllRequests()
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetNewRequests обрабатывает GET /admin/partner-requests/new, получает только новые заявки
func (h *Handler) GetNewRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.admin.GetRequestsByStatus("new")
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetPendingRequests обрабатывает GET /admin/partner-requests/pending, получает только заявки в работе
func (h *Handler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.admin.GetRequestsByStatus("pending")
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetApprovedRequests обрабатывает GET /admin/partner-requests/approved, получает только принятые заявки
func (h *Handler) GetApprovedRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.admin.GetRequestsByStatus("approved")
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetRejectedRequests обрабатывает GET /admin/partner-requests/rejected, получает только отклоненные заявки
func (h *Handler) GetRejectedRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.admin.GetRequestsByStatus("rejected")
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// ApprovePartnerRequest обрабатывает POST /admin/partner-requests/approve, принимает заявку
// автоматически создает компанию и нулевого пользователя в ней
func (h *Handler) ApprovePartnerRequest(w http.ResponseWriter, r *http.Request) {
	var req ApprovePartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.ApprovePartnerRequest(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Partner request approved",
	})
}

// RejectPartnerRequest обрабатывает POST /admin/partner-requests/reject, отклоняет заявку
func (h *Handler) RejectPartnerRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.RejectPartnerRequest(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request rejected",
		"inn":     req.ID.String(),
		"status":  "rejected",
	})
}

// GetRequest обрабатывает GET /admin/partner-requests/{id}
func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID must be UUID", http.StatusBadRequest)
		return
	}

	req, err := h.admin.GetRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// CreateAdmin обрабатывает POST /admin/create-admin, создаёт нового администратора
func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	// Получение claims из токена
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Проверка, что текущий пользователь — админ
	adminEmail, ok := claims["email"].(string)
	if !ok || adminEmail == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.admin.IsAdmin(adminEmail)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Forbidden: only admins can create new admins", http.StatusForbidden)
		return
	}

	var req struct {
		Email   string `json:"email" validate:"required,email"`
		Name    string `json:"name" validate:"required"`
		Surname string `json:"surname" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err = h.admin.CreateAdmin(req.Email, req.Name, req.Surname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Admin created successfully",
		"email":   req.Email,
	})
}

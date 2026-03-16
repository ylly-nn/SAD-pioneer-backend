package admin

import (
	"encoding/json"
	"net/http"

	"src/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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
		INN string `json:"inn" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.TakeRequestToWork(req.INN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request taken to work",
		"inn":     req.INN,
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

	err := h.admin.ApprovePartnerRequest(req.INN)
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
		INN string `json:"inn" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	err := h.admin.RejectPartnerRequest(req.INN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request rejected",
		"inn":     req.INN,
		"status":  "rejected",
	})
}

// GetRequestStatus обрабатывает GET /partner/request/{inn}, получает статус заявки
func (h *Handler) GetRequestStatus(w http.ResponseWriter, r *http.Request) {
	inn := chi.URLParam(r, "inn")
	if inn == "" {
		http.Error(w, "INN is required", http.StatusBadRequest)
		return
	}

	req, err := h.admin.GetRequestStatus(inn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

package partners

import (
	"encoding/json"
	"net/http"

	"src/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

// Handler обрабатывает HTTP-запросы для авторизации пользователей
type Handler struct {
	partners *PartnersManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(service *PartnersManager) *Handler {
	return &Handler{partners: service}
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

	err := h.partners.CreatePartnerRequest(userEmail, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Partner request created successfully",
	})
}

// GetRequestStatus обрабатывает GET /partner/request/{inn}, получает статус заявки
func (h *Handler) GetRequestStatus(w http.ResponseWriter, r *http.Request) {
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

	req, err := h.partners.GetRequestStatusByEmail(userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Handler обрабатывает HTTP-запросы для клиентов
type Handler struct {
	client *ClientManager
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(client *ClientManager) *Handler {
	return &Handler{client: client}
}

// CreateClient обрабатывает POST /client
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {

	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	client, err := h.client.CreateClient(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyEmail):
			http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrClientAlreadyExists):
			http.Error(w, "Client with this email already exists", http.StatusConflict)
		case errors.Is(err, ErrUserNotFound):
			http.Error(w, "User with this email not found in all_users", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := ClientResponse{
		ID:    client.ID,
		Email: client.Email,
		City:  client.City,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// UpdateCity обрабатывает PUT /client/city (или PATCH)
func (h *Handler) UpdateCity(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	// Получение email из claims (из тела токена)
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	var req UpdateCityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.client.UpdateCity(email, req.City); err != nil {
		switch {
		case errors.Is(err, ErrInvalidCity):
			http.Error(w, "City is not in the list of Russian cities", http.StatusBadRequest)
		case errors.Is(err, ErrEmptyEmail):
			http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrEmptyCity):
			http.Error(w, "City cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrClientNotFound):
			http.Error(w, "Client not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetCity обрабатывает GET /client/city/
func (h *Handler) GetCity(w http.ResponseWriter, r *http.Request) {

	// Извлеxtybt  данных пользователя из контекста (добавлены в middleware)
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	// Получение email из claims (из тела токена)
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	city, err := h.client.GetCityByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyEmail):
			http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrClientNotFound):
			http.Error(w, "Client not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := GetCityResponse{City: city}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

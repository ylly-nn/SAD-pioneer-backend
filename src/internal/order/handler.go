package order

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler обрабатывает HTTP-запросы для заказов
type Handler struct {
	order *OrderManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(order *OrderManager) *Handler {
	return &Handler{order: order}
}

// GetOrdersByClient обрабатывает GET /order/by-client и возвращает список всех заказов
// с добавленным email клиента из таблицы ts_users.
func (h *Handler) GetFullAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.order.GetFullAllOrders()
	if err != nil {
		log.Printf("GetOrdersByClient error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("GetOrdersByClient encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetClientOrders обрабатывает GET /order/client/{email} и возвращает полную информацию
// о заказах конкретного клиента (email из пути).
func (h *Handler) GetClientOrders(w http.ResponseWriter, r *http.Request) {

	email := chi.URLParam(r, "email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	orders, err := h.order.GetByClient(email)
	if err != nil {
		log.Printf("GetClientOrders error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("GetClientOrders encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetCompanyOrders обрабатывает GET /order/company/{inn} и возвращает полную информацию
// о заказах всех клиентов для указанной организации (по ИНН).
func (h *Handler) GetCompanyOrders(w http.ResponseWriter, r *http.Request) {

	inn := chi.URLParam(r, "inn")
	if inn == "" {
		http.Error(w, "INN is required", http.StatusBadRequest)
		return
	}

	orders, err := h.order.GetByCompany(inn)
	if err != nil {

		if errors.Is(err, ErrInnLen) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("GetCompanyOrders error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("GetCompanyOrders encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CreateOrder обрабатывает POST /order и создаёт новый заказ.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	order, err := h.order.Create(req)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("CreateOrder encode error: %v", err)
	}
}

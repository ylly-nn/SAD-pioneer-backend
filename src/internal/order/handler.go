package order

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler обрабатывает HTTP-запросы для заказов
type Handler struct {
	order *OrderManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(order *OrderManager) *Handler {
	return &Handler{order: order}
}

// GetAllOrders обрабатывает GET /order и возвращает список всех заказов.
func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.order.GetAllOrders()
	if err != nil {
		log.Printf("GetAllOrders error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("GetAllOrders encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

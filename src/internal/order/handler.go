package order

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"src/internal/timeparsing"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы для заказов
type Handler struct {
	order *OrderManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(order *OrderManager) *Handler {
	return &Handler{order: order}
}

// CreateOrder обрабатывает POST /order и создаёт новый заказ.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

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

	var req CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	order, err := h.order.Create(email, req)

	if err != nil {
		log.Printf("CreateOrdererror: %v", err)
		switch {
		case errors.Is(err, ErrBranchServIsEmpty):
			http.Error(w, ErrBranchServIsEmpty.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrStartMomentIsEmpty):
			http.Error(w, ErrStartMomentIsEmpty.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrOrderDetailsIsEmpty):
			http.Error(w, ErrOrderDetailsIsEmpty.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrTimeInPast):
			http.Error(w, ErrTimeInPast.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrTimeInFuture):
			http.Error(w, ErrTimeInFuture.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrStartMomemtNotAvailable):
			http.Error(w, ErrStartMomemtNotAvailable.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrBranchServiceNotFound):
			http.Error(w, ErrBranchServiceNotFound.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrDetailNotAvailable):
			http.Error(w, ErrDetailNotAvailable.Error(), http.StatusBadRequest)
		default:

			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("CreateOrder encode error: %v", err)
	}
}

// GetFreeTime обрабатывает GET /branch/freetime?branch_id=<uuid>&date=YYYY-MM-DD&duration=<minutes>
func (h *Handler) GetFreeTime(w http.ResponseWriter, r *http.Request) {

	branchIDStr := r.URL.Query().Get("branch_id")
	if branchIDStr == "" {
		http.Error(w, "missing branch_id", http.StatusBadRequest)
		return
	}
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		http.Error(w, "invalid branch_id", http.StatusBadRequest)
		return
	}

	//  date (формат YYYY-MM-DD)
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "missing date", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	//  duration (целое положительное число минут)
	durationStr := r.URL.Query().Get("duration")
	if durationStr == "" {
		http.Error(w, "missing duration", http.StatusBadRequest)
		return
	}
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		http.Error(w, "invalid duration, must be positive integer", http.StatusBadRequest)
		return
	}

	tzParam := r.URL.Query().Get("timezone")
	loc := time.UTC
	if tzParam != "" {
		var err error
		loc, err = timeparsing.ParseLocation(tzParam)
		if err != nil {
			http.Error(w, "invalid timezone format. Use IANA name (e.g., Europe/Moscow) or offset (e.g., +03:00)", http.StatusBadRequest)
			return
		}
	}

	slots, err := h.order.GetFreeTimeForWeek(branchID, date, duration)

	if err != nil {
		log.Printf("GetFreeTime error: %v", err)
		switch {
		case errors.Is(err, ErrBranchNotFound):
			http.Error(w, ErrBranchNotFound.Error(), http.StatusNotFound)
		case errors.Is(err, ErrDateInPast):
			http.Error(w, ErrDateInPast.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrDateInFuture):
			http.Error(w, ErrDateInFuture.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	if tzParam == "" || tzParam == "0" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(slots); err != nil {
			log.Printf("GetFreeTime encode error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := make([]DailySlotsTZ, len(slots))
	for i, day := range slots {
		// Локальная дата (полночь в заданной зоне)
		dateInLoc := time.Date(day.Date.Year(), day.Date.Month(), day.Date.Day(), 0, 0, 0, 0, loc)
		intervals := make([]time.Time, len(day.Intervals))
		for j, slot := range day.Intervals {
			intervals[j] = time.Time(slot).In(loc)
		}
		response[i] = DailySlotsTZ{
			Date:      dateInLoc,
			Intervals: intervals,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GetFreeTime encode error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// GetClientOrders обрабатывает GET /client/orders и возвращает заказы текущего аутентифицированного клиента
// Email извлекается из JWT токена, который добавляется в контекст middleware'ой AuthMiddleware.Authenticate
func (h *Handler) GetClientOrders(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	orders, err := h.order.GetByClient(email)
	if err != nil {
		if errors.Is(err, ErrOrdersNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil)
			return
		}
		log.Printf("GetClientOrders error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tzParam := r.URL.Query().Get("timezone")

	// Если tz не передан исходная структуру (с UTCTime)
	if tzParam == "" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			log.Printf("GetClientOrders encode error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	loc, err := timeparsing.ParseLocation(tzParam)
	if err != nil {
		http.Error(w, "invalid timezone format. Use IANA name (e.g., Europe/Moscow) or offset (e.g., +03:00)", http.StatusBadRequest)
		return
	}

	// Преобразуем каждый заказ в ClientOrderResponseTZ с временем в нужном поясе
	response := make([]*ClientOrderResponseTZ, len(orders))
	for i, o := range orders {
		// Преобразуем UTCTime -> time.Time и меняем зону
		start := time.Time(o.StartMoment).In(loc)
		tzOrder := &ClientOrderResponseTZ{
			ID:           o.ID,
			NameCompany:  o.NameCompany,
			City:         o.City,
			Address:      o.Address,
			Service:      o.Service,
			StartMoment:  start,
			Status:       o.Status,
			OrderDetails: o.OrderDetails,
			Sum:          o.Sum,
		}
		if o.EndMoment != nil {
			end := time.Time(*o.EndMoment).In(loc)
			tzOrder.EndMoment = &end
		}
		response[i] = tzOrder
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GetClientOrders encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetCompanyOrders обрабатывает GET /order/company/{inn} и возвращает полную информацию
// о заказах всех клиентов для указанной организации (по ИНН).
// func (h *Handler) GetCompanyOrders(w http.ResponseWriter, r *http.Request) {

// 	inn := chi.URLParam(r, "inn")
// 	if inn == "" {
// 		http.Error(w, "INN is required", http.StatusBadRequest)
// 		return
// 	}

// 	orders, err := h.order.GetByCompany(inn)
// 	if err != nil {

// 		if errors.Is(err, ErrInnLen) {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		log.Printf("GetCompanyOrders error: %v", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(orders); err != nil {
// 		log.Printf("GetCompanyOrders encode error: %v", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	}
// }

// // возвращает список всех заказов
// func (h *Handler) GetFullAllOrders(w http.ResponseWriter, r *http.Request) {
// 	orders, err := h.order.GetFullAllOrders()
// 	if err != nil {
// 		log.Printf("GetOrdersByClient error: %v", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(orders); err != nil {
// 		log.Printf("GetOrdersByClient encode error: %v", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	}
// }

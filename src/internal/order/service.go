package order

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrEmptyEmail = errors.New("email cannot be empty")
	ErrInnLen     = errors.New("inn length must be 12 characters")
)

// содержит бизнес-логику для работы с услугами.
type OrderManager struct {
	storage OrderStorage
}

// создаёт новый экземпляр OrderManager.
func NewOrderManager(storage OrderStorage) *OrderManager {
	return &OrderManager{storage: storage}
}

// возвращает список всех заказов
func (m *OrderManager) GetFullAllOrders() ([]*FullOrder, error) {
	return m.storage.GetFullAllOrders()
}

// возвращает список заказов опредеоённого киента
func (m *OrderManager) GetByClient(email string) ([]*FullOrder, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}
	return m.storage.GetByCLient(email)
}

func (m *OrderManager) GetByCompany(inn string) ([]*FullOrder, error) {
	if len(inn) != 12 {
		return nil, ErrInnLen
	}
	return m.storage.GetByCompany(inn)
}

// Create создаёт новый заказ.
func (m *OrderManager) Create(req CreateOrderRequest) (*Order, error) {
	// Проверка обязательных полей
	if req.Users == uuid.Nil {
		return nil, errors.New("users is required")
	}
	if req.ServiceByBranch == uuid.Nil {
		return nil, errors.New("service_by_branch is required")
	}
	if req.Date == nil {
		return nil, errors.New("date is required")
	}
	if req.StartTime == nil {
		return nil, errors.New("start_time is required")
	}
	if len(req.OrderDetails) == 0 {
		return nil, errors.New("order_details is required")
	}

	// Валидация структуры order_details: должен быть объектом { "название услуги": минуты }
	var details map[string]int
	if err := json.Unmarshal(req.OrderDetails, &details); err != nil {
		return nil, errors.New("order_details must be a JSON object with string keys and numeric values (minutes)")
	}
	if len(details) == 0 {
		return nil, errors.New("order_details cannot be empty")
	}
	for name, minutes := range details {
		if minutes <= 0 {
			return nil, fmt.Errorf("duration for '%s' must be positive (minutes)", name)
		}
	}

	order := Order{
		Users:           req.Users,
		ServiceByBranch: req.ServiceByBranch,
		Date:            req.Date,
		StartTime:       req.StartTime,
		OrderDetails:    req.OrderDetails,
	}
	return m.storage.Create(order)
}

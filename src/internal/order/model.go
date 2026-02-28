package order

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Соответствует таблице orders в базе данных
type Order struct {
	ID              uuid.UUID       `json:"id"`
	Users           uuid.UUID       `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	Date            *time.Time      `json:"date,omitempty"`
	StartTime       *time.Time      `json:"start_time,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// Структура заказа со всеми необходимыми данными
type FullOrder struct {
	ID              uuid.UUID       `json:"id"`
	Users           uuid.UUID       `json:"idusers"`
	Email           string          `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	InnCompany      string          `json:"inn"`
	NameCompany     string          `json:"name_company"`
	City            string          `json:"city"`
	Address         string          `json:"address"`
	Service         string          `json:"service"`
	Date            *time.Time      `json:"date,omitempty"`
	StartTime       *time.Time      `json:"start_time,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// CreateOrderRequest используется для POST /orders.
type CreateOrderRequest struct {
	Users           uuid.UUID       `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	Date            *time.Time      `json:"date,omitempty"`
	StartTime       *time.Time      `json:"start_time,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// OrderResponse \- Get order
type OrderResponse struct {
	ID              uuid.UUID       `json:"id"`
	Users           uuid.UUID       `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	Date            *time.Time      `json:"date,omitempty"`
	StartTime       *time.Time      `json:"start_time,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

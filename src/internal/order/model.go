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

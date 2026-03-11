package order

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Order представляет заказ, соответствующий таблице orders в базе данных.
type Order struct {
	ID              uuid.UUID       `json:"id"`
	Users           string          `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	StartMoment     time.Time       `json:"start_moment"`
	EndMoment       *time.Time      `json:"end_moment,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details"`
}

// BusyTime представляет временной интервал занятости (начало и конец)
type BusyTime struct {
	StartMoment time.Time `json:"start_moment"`
	EndMoment   time.Time `json:"end_moment"`
}

// DailyIntervals содержит дату и список занятых интервалов за этот день.
// Используется для возврата данных из GetBisyTimeForWeek
type DailyIntervals struct {
	Date      time.Time
	Intervals []*BusyTime
}

// DailySlots содержит дату и список времён начала доступных слотов фиксированной длительности.
// Используется для возврата данных из GetFreeTimeForWeek
type DailySlots struct {
	Date      time.Time   `json:"date"`
	Intervals []time.Time `json:"intervals"`
}

// OpenCloseBranch содержит время открытия и закрытия филиала
type OpenCloseBranch struct {
	OpenTimeBranch  time.Time
	CloseTimeBranch time.Time
}

// CreateOrderRequest используется для POST /orders.
type CreateOrderRequest struct {
	Users           string          `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	StartMoment     time.Time       `json:"start_moment"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// Структура заказа со всеми необходимыми данными
type FullOrder struct {
	ID              uuid.UUID       `json:"id"`
	Email           string          `json:"users"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch"`
	InnCompany      string          `json:"inn"`
	NameCompany     string          `json:"name_company"`
	City            string          `json:"city"`
	Address         string          `json:"address"`
	Service         string          `json:"service"`
	StartMoment     time.Time       `json:"start_moment"`
	EndMoment       *time.Time      `json:"end_moment,omitempty"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// ClientOrderResponse - упрощённая информация о заказе для клиента
type ClientOrderResponse struct {
	ID           uuid.UUID       `json:"order_id"`
	NameCompany  string          `json:"name_company"`
	City         string          `json:"city"`
	Address      string          `json:"address"`
	Service      string          `json:"service"`
	StartMoment  time.Time       `json:"start_moment"`
	EndMoment    *time.Time      `json:"end_moment,omitempty"`
	OrderDetails json.RawMessage `json:"order_details,omitempty"`
}

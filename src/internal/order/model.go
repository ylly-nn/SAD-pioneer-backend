package order

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Order представляет заказ, соответствующий таблице orders в базе данных.
type Order struct {
	ID              uuid.UUID       `json:"id" example:"e77fd339-9478-4375-82c1-215936a68b8a"`
	Users           string          `json:"users" example:"ex@mail.ru"`
	ServiceByBranch uuid.UUID       `json:"service_by_branch" example:"917e77fa-1672-4dfb-8507-d5755b31ebb3"`
	StartMoment     time.Time       `json:"start_moment" example:"2026-04-16T05:00:00Z"`
	EndMoment       *time.Time      `json:"end_moment,omitempty" example:"2026-04-16T05:20:00Z"`
	Status          string          `json:"status" example:"create"`
	OrderDetails    json.RawMessage `json:"order_details" swaggertype:"object"`
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
	Date      time.Time   `json:"date" example:"2026-03-16T00:00:00Z" format:"yyyy-mm-ddT00:00:00Z"`
	Intervals []time.Time `json:"intervals" swagertype:"array" example:"[2026-03-17T10:35:00+04:00, 2026-03-17T10:50:00+04:00, 2026-03-17T11:05:00+04:00]"`
}

// OpenCloseBranch содержит время открытия и закрытия филиала
type OpenCloseBranch struct {
	OpenTimeBranch  time.Time
	CloseTimeBranch time.Time
}

// CreateOrderRequest используется для POST /orders.
type CreateOrderRequest struct {
	ServiceByBranch uuid.UUID       `json:"service_by_branch" example:"89d74b8a-8cee-44fa-96ea-6aec1e8ad66b" format:"uuid"`
	StartMoment     time.Time       `json:"start_moment" example:"2026-03-16T11:20:00+04:00" format:"yyyy-mm-ddThh-mm-ss+hh:mm"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty" swaggertype:"object"`
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
	Status          string          `json:"status" example:"create"`
	OrderDetails    json.RawMessage `json:"order_details,omitempty"`
}

// ClientOrderResponse - упрощённая информация о заказе для клиента
type ClientOrderResponse struct {
	ID           uuid.UUID       `json:"order_id" example:"83817fd0-ffd0-478b-b1ae-b082e8581830"`
	NameCompany  string          `json:"name_company" example:"ООО \"Ромашка\""`
	City         string          `json:"city" example:"Москва"`
	Address      string          `json:"address" example:"ул. Тверская, д. 1"`
	Service      string          `json:"service" example:"Шиномонтаж"`
	StartMoment  time.Time       `json:"start_moment" example:"2026-03-16T09:30:00+04:00" format:"yyyy-mm-ddThh:mm:ss+(Z)hh:mm"`
	EndMoment    *time.Time      `json:"end_moment" example:"2026-03-16T11:05:00+04:00" format:"yyyy-mm-ddThh:mm:ss+(Z)hh:mm"`
	Status       string          `json:"status" example:"create"`
	OrderDetails json.RawMessage `json:"order_details" swaggertype:"object"`
}

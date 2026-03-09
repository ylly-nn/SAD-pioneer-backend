package branch

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Branch соответствует таблице branch - для внутреннейй передачи данных
type Branch struct {
	ID        uuid.UUID `json:"branch_id"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Inn       string    `json:"inn_company"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

// Используется для обработки POST /branch
type CreateBranchRequest struct {
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Inn       string    `json:"inn_company"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

// BranchServ представляет связь филиала и услуги с деталями услуги и временем занятости.
type BranchServ struct {
	ID             uuid.UUID       `json:"id"`
	Branch         uuid.UUID       `json:"branch"`
	Service        uuid.UUID       `json:"service"`
	ServiceDetails json.RawMessage `json:"service_detalis"`
	BusyTime       json.RawMessage `json:"busy_time,omitempty"`
}

// CreateBranchServRequest содержит данные для создания записи branch_services.
type CreateBranchServRequest struct {
	ID             uuid.UUID       `json:"id"`
	Branch         uuid.UUID       `json:"branch"`
	Service        uuid.UUID       `json:"service"`
	ServiceDetails json.RawMessage `json:"service_detalis"`
}

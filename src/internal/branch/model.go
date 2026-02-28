package branch

import (
	"encoding/json"

	"github.com/google/uuid"
)

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

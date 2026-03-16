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

// Использзуется для Get /branch?city=<city>&service=<service>
type BrancByCityServ struct {
	BranchServId uuid.UUID `json:"id_branchserv" example:"0bb7a20d-4ffc-46cd-b5e5-a549a179ce2a" format:"uuid"`
	BranchId     uuid.UUID `json:"id_branch" example:"0bb7a20d-4ffc-46cd-b5e5-a549a179ce2a" format:"uuid"`
	Address      string    `json:"address" example:"ул. Тверская, д. 1"`
	CompanyName  string    `json:"org_short_name" example:"ООО Ромашка"`
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

// Список уточнений услуги - деталь длительность
type ServiceDetails struct {
	Detail   string `json:"detail" example:"Мойка колёс"`
	Duration int    `json:"duration_min" example:"35"`
}

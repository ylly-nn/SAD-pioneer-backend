package company

import (
	"encoding/json"

	"github.com/google/uuid"

	"src/internal/timeparsing"
)

// Соответсвует таблице Company
// Используется для внутренней передачи данных
type Company struct {
	INN          string  `json:"inn"`
	KPP          *string `json:"kpp,omitempty"`
	OGRN         *string `json:"ogrn,omitempty"`
	OrgName      *string `json:"org_name,omitempty"`
	OrgShortName *string `json:"org_short_name,omitempty"`
}

// CreateCompanyRequest содержит данные для создания компании через POST /company
type CreateCompanyRequest struct {
	Company
}

// CompanyResponse содержит данные для ответа на запросы, связанные с компаниями
type CompanyResponse struct {
	INN          string  `json:"inn"`
	KPP          *string `json:"kpp,omitempty"`
	OGRN         *string `json:"ogrn,omitempty"`
	OrgName      *string `json:"org_name,omitempty"`
	OrgShortName *string `json:"org_short_name,omitempty"`
}

// IsPartnerUsers используется для проверки есть ли у пользователся организация
type IsPartnersUsers struct {
	IsPartner bool
	Inn       string
}

// PartnersUsers используется для передачи email и inn - если есть
type PartnersUsers struct {
	Email string
	Inn   string
}

// Соответсвует таблице branches, используется в get /company/branches
type CompanyBranch struct {
	ID        uuid.UUID            `json:"branch_id"`
	City      string               `json:"city"`
	Address   string               `json:"address"`
	Inn       string               `json:"inn_company"`
	OpenTime  timeparsing.TimeOnly `json:"open_time"`
	CloseTime timeparsing.TimeOnly `json:"close_time"`
}

// Филиал с соответсвуюими ему услугами
type CompanyBranchWithServ struct {
	City      string               `json:"city"`
	Address   string               `json:"address"`
	OpenTime  timeparsing.TimeOnly `json:"open_time"`
	СloseTime timeparsing.TimeOnly `json:"close_time"`
	Services  []*ServiceInBranch   `json:"services"`
}

// Cтрутура описывающая услугу для CompanyBranchWithServ struct
type ServiceInBranch struct {
	BranchServId uuid.UUID `json:"branch_serv_id"`
	ServiceId    uuid.UUID `json:"service_id"`
	ServiceName  string    `json:"service_name"`
}

// Струтура для ответа на Get /company/branch/service/{branchserID}
type CompanyServDetailsResponse struct {
	Detail   string `json:"detail"`
	Duration int    `json:"duration_min"`
}

// соответствует таблице branch_serv
type BranchServ struct {
	ID             uuid.UUID       `json:"id"`
	Branch         uuid.UUID       `json:"branch"`
	Service        uuid.UUID       `json:"service"`
	ServiceDetails json.RawMessage `json:"service_detalis"`
}

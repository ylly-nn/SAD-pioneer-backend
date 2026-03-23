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
	INN          string  `json:"inn" example:"234567890123"`
	KPP          *string `json:"kpp,omitempty" example:"234567891"`
	OGRN         *string `json:"ogrn,omitempty" example:"2345678901234"`
	OrgName      *string `json:"org_name,omitempty" example:"АО Технопром"`
	OrgShortName *string `json:"org_short_name,omitempty" example:"Технопром"`
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
	ID        uuid.UUID            `json:"branch_id" example:"9eebb3b9-5b35-4007-9d4f-2f4141786b45" `
	City      string               `json:"city" example:"Москва" `
	Address   string               `json:"address" example:"ул. Тверская, 1" `
	Inn       string               `json:"inn_company" example:"123456789012" `
	OpenTime  timeparsing.TimeOnly `json:"open_time" example:"10:00:00+00:00" format:"hh:mm:ss+hh:mm"`
	CloseTime timeparsing.TimeOnly `json:"close_time" example:"18:00:00+00:00" format:"hh:mm:ss+hh:mm"`
}

// Филиал с соответсвуюими ему услугами
type CompanyBranchWithServ struct {
	City      string               `json:"city" example:"Санкт-Петербург"`
	Address   string               `json:"address" example:"Невский пр., 10"`
	OpenTime  timeparsing.TimeOnly `json:"open_time" example:"10:00:00+00:00"  format:"hh:mm:ss+hh:mm"`
	СloseTime timeparsing.TimeOnly `json:"close_time" example:"18:00:00+00:00"  format:"hh:mm:ss+hh:mm"`
	Services  []*ServiceInBranch   `json:"services"`
}

// Cтрутура описывающая услугу для CompanyBranchWithServ struct
type ServiceInBranch struct {
	BranchServId uuid.UUID `json:"branch_serv_id" example:"6fdd2352-ffc4-4140-b54c-67657f841c1c" `
	ServiceId    uuid.UUID `json:"service_id" example:"03db1f58-2bbd-481c-8d93-b2828871b376" `
	ServiceName  string    `json:"service_name" example:"мойка" `
}

// Струтура для ответа на Get /company/branch/service/{branchserID}
type CompanyServDetailsResponse struct {
	Detail   string `json:"detail" example:"Мойка салона"`
	Duration int    `json:"duration_min" example:"40"`
}

// соответствует таблице branch_serv
type BranchServ struct {
	ID             uuid.UUID       `json:"id"`
	Branch         uuid.UUID       `json:"branch"`
	Service        uuid.UUID       `json:"service"`
	ServiceDetails json.RawMessage `json:"service_detalis"`
}

// AddUserRequest - запрос на добавление нового пользователя
type AddUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// AddBranchRequest - запрос на добавление нового филиала
type AddBranchRequest struct {
	City      string               `json:"city" validate:"required"`
	Address   string               `json:"address" validate:"required"`
	OpenTime  timeparsing.TimeOnly `json:"open_time" validate:"required"`
	CloseTime timeparsing.TimeOnly `json:"close_time" validate:"required"`
}

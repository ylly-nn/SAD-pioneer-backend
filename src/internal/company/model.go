package company

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"src/internal/timeparsing"
)

type OrderStatus string

const (
	OrderStatusCreate  OrderStatus = "create"  // заказ создан
	OrderStatusApprove OrderStatus = "approve" // заказ подтверждён
	OrderStatusReject  OrderStatus = "reject"  // заказ отклонён
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
type ServDetails struct {
	Detail   string `json:"detail" example:"Мойка салона"`
	Duration int    `json:"duration_min" example:"40"`
}

type ServPrice struct {
	Detail string  `json:"detail" example:"Мойка салона"`
	Price  float32 `json:"price" example:"560.12"`
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
	City      string               `json:"city" example:"Москва" validate:"required"`
	Address   string               `json:"address" example:"Улица тверская дом 1" validate:"required"`
	OpenTime  timeparsing.TimeOnly `json:"open_time" example:"09:00:00+03:00" validate:"required"`
	CloseTime timeparsing.TimeOnly `json:"close_time" example:"17:00:00+03:00" validate:"required"`
}

// AddServiceToBranch - запрос на добавление новой услуги в филиал
type AddServiceToBranch struct {
	BranchID  uuid.UUID `json:"branch_id"  example:"9eebb3b9-5b35-4007-9d4f-2f4141786b45" validate:"required"`
	ServiceID uuid.UUID `json:"service_id" example:"03db1f58-2bbd-481c-8d93-b2828871b376" validate:"required"`
}

type CompanyBranchOrderResponse struct {
	BranchID uuid.UUID `json:"branch_id"`
	City     string    `json:"city"`
	Address  string    `json:"address"`
	Orders   []*CompanyOrder
}

type CompanyOrder struct {
	ID              uuid.UUID   `json:"id" example:"e77fd339-9478-4375-82c1-215936a68b8a"`
	Users           string      `json:"users" example:"ex@mail.ru"`
	ServiceByBranch uuid.UUID   `json:"service_by_branch" example:"917e77fa-1672-4dfb-8507-d5755b31ebb3"`
	NameService     string      `json:"name_service" example:"автомойка"`
	StartMoment     time.Time   `json:"start_moment" example:"2026-04-16T05:00:00Z"`
	EndMoment       *time.Time  `json:"end_moment,omitempty" example:"2026-04-16T05:20:00Z"`
	Status          OrderStatus `json:"status" example:"create"`
	OrderDetails    []ServDetails
}

// Запрос для обработчика для добавления детали
type AddServDetailRequest struct {
	BranchServID uuid.UUID `json:"branchserv_id" validate:"required" example:"6fdd2352-ffc4-4140-b54c-67657f841c1c"`
	Detail       string    `json:"detail" validate:"required,min=3,max=255" example:"Мойка салона"`
	Duration     int       `json:"duration" validate:"required,min=1,max=1439" example:"40"`
	Price        *float32  `json:"price" example:"700.50" validate:"required,gt=0,lt=1000000000000"`
}

type Service struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ServUpdateResponse struct {
	Detail   string  `json:"detail" example:"Мойка салона"`
	Duration int     `json:"duration_min" example:"40"`
	Price    float32 `json:"price" example:"560.12"`
}

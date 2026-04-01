package partners

import "time"

// PartnerRequest - заявка на регистрацию организации
type PartnerRequest struct {
	Status    string `json:"status" db:"status"`
	UserEmail string `json:"user_email" db:"user_email"`

	INN          string `json:"inn" db:"inn" validate:"required,inn"`
	KPP          string `json:"kpp" db:"kpp" validate:"required,kpp"`
	OGRN         string `json:"ogrn" db:"ogrn" validate:"required,ogrn"`
	OrgName      string `json:"org_name" db:"org_name" validate:"required,org_name"`
	OrgShortName string `json:"org_short_name" db:"org_short_name" validate:"required,org_short_name"`

	Name       string `json:"name" db:"name" validate:"required,person_name"`
	Surname    string `json:"surname" db:"surname" validate:"required,person_name"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic" validate:"required,person_name"`
	Email      string `json:"email" db:"email" validate:"required,email"`
	Phone      string `json:"phone" db:"phone_number" validate:"required,phone"`
	Info       string `json:"info,omitempty" db:"info" validate:"omitempty"`

	CreatedAt time.Time  `json:"created_at"  example:"2026-03-30T06:06:47.181805Z" db:"created_at"`
	LastUsed  *time.Time `json:"last_used,omitempty" example:"2026-03-30T06:07:27.657019Z" db:"last_used"`
}

// Company - данные организации
type Company struct {
	INN          string `json:"inn" db:"inn"`
	KPP          string `json:"kpp" db:"kpp"`
	OGRN         string `json:"ogrn" db:"ogrn"`
	OrgName      string `json:"org_name" db:"org_name"`
	OrgShortName string `json:"org_short_name" db:"org_short_name"`
}

// PartnerRequestRequest - запрос на создание заявки партнера
type PartnerRequestRequest struct {
	Status    string `json:"status" db:"status"`
	UserEmail string `json:"-" db:"user_email"`

	INN          string `json:"inn" db:"inn" example:"123456789012" validate:"required,inn"`
	KPP          string `json:"kpp" db:"kpp" example:"123456789" validate:"required,kpp"`
	OGRN         string `json:"ogrn" db:"ogrn" example:"1234567890123"  validate:"required,ogrn"`
	OrgName      string `json:"org_name" db:"org_name" example:"ООО Ромашка" validate:"required,org_name"`
	OrgShortName string `json:"org_short_name" db:"org_short_name" example:"Ромашка" validate:"required,org_short_name"`

	Name       string `json:"name" db:"name" example:"Иван" validate:"required,person_name"`
	Surname    string `json:"surname" db:"surname" example:"Иванов" validate:"required,person_name"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic" example:"Иванович" validate:"required,person_name"`
	Email      string `json:"email" db:"email" example:"ivan@example.com" validate:"required,email"`
	Phone      string `json:"phone" db:"phone_number" example:"9871111111" validate:"required,phone"`
	Info       string `json:"info,omitempty" db:"info" example:"Дополнительная информация" validate:"omitempty"`
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

// Ошибки
var (
	ErrUserNotFound        = "user not found"
	ErrInvalidPassword     = "invalid password"
	ErrUserAlreadyExists   = "user already exists"
	ErrInvalidCode         = "invalid verification code"
	ErrCodeExpired         = "verification code expired"
	ErrInvalidRefreshToken = "invalid refresh token"
	ErrRefreshTokenExpired = "refresh token expired"
)

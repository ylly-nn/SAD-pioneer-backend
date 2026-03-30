package partners

import "time"

// PartnerRequest - заявка на регистрацию организации
type PartnerRequest struct {
	Status    string `json:"status" db:"status"`
	UserEmail string `json:"user_email" db:"user_email"`

	INN          string `json:"inn" db:"inn"`
	KPP          string `json:"kpp" db:"kpp"`
	OGRN         string `json:"ogrn" db:"ogrn"`
	OrgName      string `json:"org_name" db:"org_name"`
	OrgShortName string `json:"org_short_name" db:"org_short_name"`

	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phone" db:"phone_number"`
	Info       string `json:"info,omitempty" db:"info"`

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

	INN          string `json:"inn" db:"inn" example:"123456789012"`
	KPP          string `json:"kpp" db:"kpp" example:"123456789"`
	OGRN         string `json:"ogrn" db:"ogrn" example:"1234567890123"`
	OrgName      string `json:"org_name" db:"org_name" example:"ООО Ромашка"`
	OrgShortName string `json:"org_short_name" db:"org_short_name" example:"Ромашка"`

	Name       string `json:"name" db:"name" example:"Иван"`
	Surname    string `json:"surname" db:"surname" example:"Иванов"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic" example:"Иванович"`
	Email      string `json:"email" db:"email" example:"ivan@example.com"`
	Phone      string `json:"phone" db:"phone_number" example:"9871111111"`
	Info       string `json:"info,omitempty" db:"info" example:"Дополнительная информация"`
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

package admin

// PartnerRequest - заявка на регистрацию организации
type PartnerRequest struct {
	Status    string `json:"status" db:"status"`
	UserEmail string `json:"-" db:"user_email"`

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
	INN          string `json:"inn" validate:"required"`
	KPP          string `json:"kpp" validate:"required"`
	OGRN         string `json:"ogrn" validate:"required"`
	OrgName      string `json:"org_name" validate:"required"`
	OrgShortName string `json:"org_short_name" validate:"required"`

	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	Info       string `json:"info"`
}

// ApprovePartnerRequest - запрос на одобрение заявки
type ApprovePartnerRequest struct {
	INN string `json:"inn" validate:"required"`
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

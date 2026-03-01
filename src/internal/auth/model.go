package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Роль
const (
	UserTypeClient = "client"
)

// Соответствует таблице all_users в базе данных
type User struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"-" db:"password"`
	TypeUser string `json:"type_user" db:"type_user"`
}

// Данные регистрации
type RegistrationData struct {
	Email     string
	Password  string
	Code      string
	ExpiresAt time.Time
}

// Данные refresh токена
type RefreshTokenData struct {
	Token      string
	Email      string
	ExpiresAt  time.Time
	LastUsed   time.Time
	DeviceInfo string
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=20,password"`
}

// Запрос на подтверждение кода
type VerifyCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

// Запрос на логин
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Ответ с токенами
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Данные, которые хранятся внутри jwt токена
type Claims struct {
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
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

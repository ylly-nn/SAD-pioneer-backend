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
	Email    string `json:"email" example:"example@gmail.com" validate:"required,email,min=6,max=64"`
	Password string `json:"password" example:"1A_password" validate:"required,min=8,max=24,password"`
}

// Запрос на подтверждение кода
type VerifyCodeRequest struct {
	Email string `json:"email" example:"example@gmail.com" validate:"required,email"`
	Code  string `json:"code" example:"000000" validate:"required,len=6"`
}

// Запрос на логин
type LoginRequest struct {
	Email    string `json:"email" example:"example@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"1A_password" validate:"required"`
}

// Запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." validate:"required"`
}

// Ответ с токенами
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"900"`
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

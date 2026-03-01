package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	configPkg "src/internal/config"
)

type Config struct {
	JWTSecretKey    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	VerificationTTL time.Duration
}

// Содержит бизнес-логику для работы с авторизацией
type AuthManager struct {
	userStorage         UserStorage
	refreshTokenStorage RefreshTokenStorage
	verificationStorage VerificationStorage
	emailSender         configPkg.EmailSender
	config              Config
}

// Создает новый экземпляр сервиса
func NewAuthManager(
	userStorage UserStorage,
	refreshTokenStorage RefreshTokenStorage,
	verificationStorage VerificationStorage,
	emailSender configPkg.EmailSender,
	config Config,
) *AuthManager {
	return &AuthManager{
		userStorage:         userStorage,
		refreshTokenStorage: refreshTokenStorage,
		verificationStorage: verificationStorage,
		emailSender:         emailSender,
		config:              config,
	}
}

// Регистрация, отправка кода подтверждения
func (s *AuthManager) Register(email, password string) error {
	// Проверка существования пользователя
	existingUser, err := s.userStorage.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf(ErrUserAlreadyExists)
	}

	// Хеширование пароля
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Генерация кода
	code := s.generateVerificationCode()

	// Сохранение данных
	data := &RegistrationData{
		Email:     email,
		Password:  hashedPassword,
		Code:      code,
		ExpiresAt: time.Now().Add(s.config.VerificationTTL),
	}

	if err := s.verificationStorage.Save(data); err != nil {
		return fmt.Errorf("failed to save verification data: %w", err)
	}

	// Отправление кода
	if err := s.emailSender.SendVerificationCode(email, code); err != nil {
		s.verificationStorage.Delete(email)
		return fmt.Errorf("failed to send verification code: %w", err)
	}

	return nil
}

// Верификация кода, окончательная регистрация (добавление пользователя в БД)
func (s *AuthManager) VerifyCode(email, code string) error {
	// Получение данных из БД
	data, err := s.verificationStorage.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get verification data: %w", err)
	}
	if data == nil {
		return fmt.Errorf(ErrInvalidCode)
	}

	// Проверка срока действия кода
	if time.Now().After(data.ExpiresAt) {
		s.verificationStorage.Delete(email)
		return fmt.Errorf(ErrCodeExpired)
	}

	// Проверка кода
	if data.Code != code {
		return fmt.Errorf(ErrInvalidCode)
	}

	// Создание пользователя
	user := &User{
		Login:    email,
		Password: data.Password,
		TypeUser: UserTypeClient,
	}

	if err := s.userStorage.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.verificationStorage.Delete(email)

	return nil
}

// Вход в систему
func (s *AuthManager) Login(email, password string) (*TokenResponse, error) {
	// Получение пользователя
	user, err := s.userStorage.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf(ErrUserNotFound)
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf(ErrInvalidPassword)
	}

	// Генерирация токенов
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// Обновление токенов
func (s *AuthManager) RefreshTokens(refreshToken string) (*TokenResponse, error) {
	tokenData, err := s.refreshTokenStorage.GetByToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	if tokenData == nil {
		return nil, fmt.Errorf(ErrInvalidRefreshToken)
	}

	// Получение срока действия
	if time.Now().After(tokenData.ExpiresAt) {
		s.refreshTokenStorage.Delete(refreshToken)
		return nil, fmt.Errorf(ErrRefreshTokenExpired)
	}

	// Получения пользователя
	user, err := s.userStorage.GetByEmail(tokenData.Email)
	if err != nil || user == nil {
		s.refreshTokenStorage.Delete(refreshToken)
		return nil, fmt.Errorf(ErrUserNotFound)
	}

	// Удаление старого токена
	s.refreshTokenStorage.Delete(refreshToken)

	// Генерация нового
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// Выход из системы
func (s *AuthManager) Logout(refreshToken string) error {
	return s.refreshTokenStorage.Delete(refreshToken)
}

// Создание пары токенов
func (s *AuthManager) generateTokenPair(user *User) (*TokenResponse, error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"email":     user.Login,
		"user_type": user.TypeUser,
		"type":      "access",
		"exp":       time.Now().Add(s.config.AccessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"email": user.Login,
		"type":  "refresh",
		"exp":   time.Now().Add(s.config.RefreshTokenTTL).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	// Сохранение refresh token в БД
	refreshData := &RefreshTokenData{
		Token:     refreshTokenString,
		Email:     user.Login,
		ExpiresAt: time.Now().Add(s.config.RefreshTokenTTL),
		LastUsed:  time.Now(),
	}

	if err := s.refreshTokenStorage.Save(refreshData); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.config.AccessTokenTTL.Seconds()),
	}, nil
}

// Генерация кода
func (s *AuthManager) generateVerificationCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)[:6]
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

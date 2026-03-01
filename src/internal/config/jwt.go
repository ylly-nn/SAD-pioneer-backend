package config

import (
	"fmt"
	"os"
	"time"
)

// Настройки для JWT
type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	VerificationTTL time.Duration
}

// Загрузка конфигурации JWT из env
func LoadJWTConfig() (*JWTConfig, error) {
	secretKey := os.Getenv("JWT_SECRET")

	accessTTL, err := parseDurationEnv("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := parseDurationEnv("REFRESH_TOKEN_TTL", 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %w", err)
	}

	verificationTTL, err := parseDurationEnv("VERIFICATION_TTL", 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("invalid VERIFICATION_TTL: %w", err)
	}

	return &JWTConfig{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		VerificationTTL: verificationTTL,
	}, nil
}

// Функция для корректного переноса длительности из env
func parseDurationEnv(key string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(value)
}

package auth

import (
	"database/sql"
	"errors"
	"sync"
	"time"
)

// Интерфейс для работы с пользователями
type UserStorage interface {
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	UpdatePassword(email, password string) error
}

// Интерфейс для работы с refresh токенами
type RefreshTokenStorage interface {
	Save(token *RefreshTokenData) error
	GetByToken(token string) (*RefreshTokenData, error)
	Delete(token string) error
	DeleteAllByEmail(email string) error
}

// Интерфейс для работы с кодами подтверждения
type VerificationStorage interface {
	Save(data *RegistrationData) error
	GetByEmail(email string) (*RegistrationData, error)
	Delete(email string) error
}

// PostgresUserStorage реализация для PostgreSQL
type PostgresUserStorage struct {
	db *sql.DB
}

func NewPostgresUserStorage(db *sql.DB) *PostgresUserStorage {
	return &PostgresUserStorage{db: db}
}

func (s *PostgresUserStorage) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT login, password, type_user FROM all_users WHERE login = $1`

	err := s.db.QueryRow(query, email).Scan(
		&user.Login,
		&user.Password,
		&user.TypeUser,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgresUserStorage) Create(user *User) error {
	query := `INSERT INTO all_users (login, password, type_user) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query,
		user.Login,
		user.Password,
		user.TypeUser,
	)

	return err
}

func (s *PostgresUserStorage) UpdatePassword(email, password string) error {
	query := `UPDATE all_users SET password = $1 WHERE login = $2`

	_, err := s.db.Exec(query, password, email)
	return err
}

// PostgresRefreshTokenStorage реализация для refresh токенов
type PostgresRefreshTokenStorage struct {
	db *sql.DB
}

func NewPostgresRefreshTokenStorage(db *sql.DB) *PostgresRefreshTokenStorage {
	storage := &PostgresRefreshTokenStorage{db: db}
	go storage.cleanupLoop()
	return storage
}

func (s *PostgresRefreshTokenStorage) Save(token *RefreshTokenData) error {
	query := `INSERT INTO refresh_tokens (token, email, expires_at, last_used) 
              VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query,
		token.Token,
		token.Email,
		token.ExpiresAt,
		time.Now(),
	)

	return err
}

func (s *PostgresRefreshTokenStorage) GetByToken(token string) (*RefreshTokenData, error) {
	var data RefreshTokenData
	query := `UPDATE refresh_tokens 
              SET last_used = CURRENT_TIMESTAMP 
              WHERE token = $1 
              RETURNING token, email, expires_at, last_used`

	err := s.db.QueryRow(query, token).Scan(&data.Token, &data.Email, &data.ExpiresAt, &data.LastUsed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (s *PostgresRefreshTokenStorage) Delete(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := s.db.Exec(query, token)
	return err
}

func (s *PostgresRefreshTokenStorage) DeleteAllByEmail(email string) error {
	query := `DELETE FROM refresh_tokens WHERE email = $1`
	_, err := s.db.Exec(query, email)
	return err
}

// Запуск очистки просроченных токенов каждый час
func (s *PostgresRefreshTokenStorage) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		_ = s.CleanExpired()
	}
}

// Удаление просроченных токенов
func (s *PostgresRefreshTokenStorage) CleanExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := s.db.Exec(query)
	return err
}

// Хранение кодов верификации в памяти
type MemoryVerificationStorage struct {
	mu       sync.RWMutex
	codes    map[string]*RegistrationData
	attempts map[string][]time.Time
}

func NewMemoryVerificationStorage() *MemoryVerificationStorage {
	storage := &MemoryVerificationStorage{
		codes:    make(map[string]*RegistrationData),
		attempts: make(map[string][]time.Time),
	}
	// Запуск очистки просроченных кодов
	go storage.cleanupLoop()
	return storage
}

// Проверка на возможность зарегистрироваться
func (s *MemoryVerificationStorage) CanRegister(email string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Получение список попыток для этого email
	attempts, exists := s.attempts[email]
	if !exists {
		s.attempts[email] = []time.Time{now}
		return true
	}

	// Удаление попыток старше 1 минуты
	var recent []time.Time
	for _, t := range attempts {
		if now.Sub(t) < time.Minute {
			recent = append(recent, t)
		}
	}

	if len(recent) >= 3 { // максимум 3 попытки в минуту
		return false
	}

	// Добавление текущей попытки
	recent = append(recent, now)
	s.attempts[email] = recent
	return true
}

// Сохранение кодов верификации
func (s *MemoryVerificationStorage) Save(data *RegistrationData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.codes[data.Email] = data
	return nil
}

// Получение данных верификации по email
func (s *MemoryVerificationStorage) GetByEmail(email string) (*RegistrationData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.codes[email]
	if !exists {
		return nil, nil
	}

	// Проверка на время действия кода
	if time.Now().After(data.ExpiresAt) {
		return nil, nil
	}

	return data, nil
}

// Удаление данных верификации
func (s *MemoryVerificationStorage) Delete(email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.codes, email)
	return nil
}

// Удаление просроченных кодов
func (s *MemoryVerificationStorage) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for email, attempts := range s.attempts {
		var recent []time.Time
		for _, t := range attempts {
			if now.Sub(t) < time.Minute {
				recent = append(recent, t)
			}
		}
		if len(recent) == 0 {
			delete(s.attempts, email)
		} else {
			s.attempts[email] = recent
		}
	}
}

// Запуск очистки просроченных кодов каждую минуту
func (s *MemoryVerificationStorage) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

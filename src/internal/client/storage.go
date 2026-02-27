package client

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"src/internal/db"
)

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrClientNotFound      = errors.New("client not found")
	ErrUserNotFound        = errors.New("user with this email not found in all_users")
)

// ClientStorage определяет методы для работы c владельцами тс в базе данных
type ClientStorage interface {
	Create(email string) (*Client, error)

	UpdateCity(email string, city string) error

	GetCityByEmail(email string) (string, error)

	//Delete - делать из all_users
	//Delete(email string)
}

// PostgresClientStorage реализует ClientStorage для PostgreSQL.
type PostgresClientStorage struct {
	*db.Storage
}

// NewPostgresClientStorage создаёт новый экземпляр PostgresClientStorage.
func NewPostgresClientStorage(sqlDB *sql.DB) *PostgresClientStorage {
	return &PostgresClientStorage{Storage: db.NewStorage(sqlDB)}
}

// Create - создаёт нового клиента если email есть в all_users
func (s *PostgresClientStorage) Create(email string) (*Client, error) {
	client := &Client{
		ID:    uuid.New(),
		Email: email,
		City:  nil,
	}

	_, err := s.DB.Exec(
		`INSERT INTO ts_users (email) VALUES ($1)`,
		client.Email,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return nil, ErrClientAlreadyExists
			case "23503": // foreign_key_violation
				return nil, ErrUserNotFound
			}
		}
		return nil, err
	}
	return client, nil
}

// UpdateCity обновляет город клиента по email.
// Если клиент с указанным email не найден, возвращает ErrClientNotFound.
func (s *PostgresClientStorage) UpdateCity(email string, city string) error {
	result, err := s.DB.Exec(
		`UPDATE ts_users SET city = $2 WHERE email = $1`,
		email, city,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}

// GetCityByEmail возвращает город клиента по email.
// Если клиент не найден, возвращает ErrClientNotFound.
// Если город не задан (NULL), возвращается пустая строка.
func (s *PostgresClientStorage) GetCityByEmail(email string) (string, error) {
	var city sql.NullString
	err := s.DB.QueryRow(
		`SELECT city FROM ts_users WHERE email = $1`,
		email,
	).Scan(&city)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrClientNotFound
		}
		return "", err
	}

	// Если город NULL, возвращаем пустую строку
	if !city.Valid {
		return "", nil
	}
	return city.String, nil
}

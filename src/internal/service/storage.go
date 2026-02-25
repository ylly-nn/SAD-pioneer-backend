package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"src/internal/db"
)

// Ошибки которые может возвращать service/storage
var (
	// ErrServiceNotFound возвращается, когда услуга с указанным ID не найдена.
	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceAlreadyExists = errors.New("service already exists")
)

// ServiceStorage определяет методы для работы с услугами в базе данных.
type ServiceStorage interface {
	GetAll() ([]Service, error)

	Create(name string) (*Service, error)

	Delete(id uuid.UUID) error
}

// PostgresServiceStorage реализует ServiceStorage для PostgreSQL.
type PostgresServiceStorage struct {
	*db.Storage
}

// NewPostgresServiceStorage создаёт новый экземпляр PostgresServiceStorage.
func NewPostgresServiceStorage(sqlDB *sql.DB) *PostgresServiceStorage {
	return &PostgresServiceStorage{Storage: db.NewStorage(sqlDB)}
}

// GetAll возвращает список всех услуг.
func (s *PostgresServiceStorage) GetAll() ([]Service, error) {
	rows, err := s.DB.Query(`SELECT id, name FROM services ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query all services: %w", err)
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var svc Service
		if err := rows.Scan(&svc.ID, &svc.Name); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}
		services = append(services, svc)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return services, nil
}

// Create - используется для создания новой услуги
func (s *PostgresServiceStorage) Create(name string) (*Service, error) {
	var service Service
	err := s.DB.QueryRow(`
        INSERT INTO services (name) 
        VALUES ($1) 
        RETURNING id, name
    `, name).Scan(&service.ID, &service.Name)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("create service: %w", ErrServiceAlreadyExists)
		}
		return nil, fmt.Errorf("create service: %w", err)
	}
	return &service, nil
}

// // Delete удаляет услугу по идентификатору.
// // Если услуга не найдена, возвращает ошибку ErrServiceNotFound.
func (s *PostgresServiceStorage) Delete(id uuid.UUID) error {
	result, err := s.DB.Exec(`DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete service %v: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for delete %v: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete service %v: %w", id, ErrServiceNotFound)
	}
	return nil
}

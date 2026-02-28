package branch

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"src/internal/db"
)

// ErrBranchServiceExists возвращается, когда запись с таким branch и service уже существует.
var ErrBranchServiceExists = errors.New("branch and service already exists")

// BranchStorage определяет методы для работы с хранилищем branch_services.
type BranchStorage interface {
	Create(branchServ BranchServ) (*BranchServ, error)
}

// PostgresBranchStorage реализует BranchStorage для PostgreSQL.
type PostgresBranchStorage struct {
	*db.Storage
}

// NewPostgresBranchStorage создаёт новый экземпляр PostgresBranchStorage.
func NewPostgresBranchStorage(sqlDB *sql.DB) *PostgresBranchStorage {
	return &PostgresBranchStorage{Storage: db.NewStorage(sqlDB)}
}

// Create добавляет новую запись в таблицу branch_services.
// Все поля, кроме BusyTime, обязательны. ID генерируется базой.
// Возвращает ошибку ErrBranchServiceExists, если пара branch+service уже существует.
// TODO(ylly): Вынести обрпботку ошибки  ErrBranchServiceExists в servise - через get
func (s *PostgresBranchStorage) Create(bs BranchServ) (*BranchServ, error) {
	if bs.Branch == uuid.Nil {
		return nil, fmt.Errorf("branch is required")
	}
	if bs.Service == uuid.Nil {
		return nil, fmt.Errorf("service is required")
	}
	if bs.ServiceDetails == nil {
		return nil, fmt.Errorf("service_detalis is required and must not be nil")
	}

	// Cуществует ли уже запись с таким branch и service
	var exists bool
	err := s.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM branch_services 
			WHERE branch = $1 AND service = $2
		)
	`, bs.Branch, bs.Service).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return nil, ErrBranchServiceExists
	}

	var id uuid.UUID
	err = s.DB.QueryRow(`
		INSERT INTO branch_services (branch, service, service_detalis, busy_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, bs.Branch, bs.Service, bs.ServiceDetails, bs.BusyTime).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert branch_service: %w", err)
	}

	created := &BranchServ{
		ID:             id,
		Branch:         bs.Branch,
		Service:        bs.Service,
		ServiceDetails: bs.ServiceDetails,
		BusyTime:       bs.BusyTime,
	}
	return created, nil
}

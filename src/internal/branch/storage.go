package branch

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"src/internal/db"
)

// ErrBranchServiceExists возвращается, когда запись с таким branch и service уже существует.
var (
	ErrBranchServiceExists   = errors.New("branch and service already exists")
	ErrBranchServiceNotFound = errors.New("branch service not found")
)

// BranchStorage определяет методы для работы с хранилищем branch_services.
type BranchStorage interface {
	CreateBranchServ(branchServ BranchServ) (*BranchServ, error)

	CreateBranch(branch Branch) (*Branch, error)

	GetBranchByCityServ(city string, serviceID string) ([]*BrancByCityServ, error)

	GetServiceDetails(branchServID uuid.UUID) ([]*ServiceDetails, []*ServPrice, error)
}

// PostgresBranchStorage реализует BranchStorage для PostgreSQL.
type PostgresBranchStorage struct {
	*db.Storage
}

// NewPostgresBranchStorage создаёт новый экземпляр Postg(resBranchStorage.
func NewPostgresBranchStorage(sqlDB *sql.DB) *PostgresBranchStorage {
	return &PostgresBranchStorage{Storage: db.NewStorage(sqlDB)}
}

// Получение деталей услуги по branch_serv
func (s *PostgresBranchStorage) GetServiceDetails(branchServID uuid.UUID) ([]*ServiceDetails, []*ServPrice, error) {
	var detailsStr sql.NullString
	var priceStr sql.NullString

	err := s.DB.QueryRow(`SELECT service_detalis, price FROM branch_services WHERE id = $1`, branchServID).Scan(&detailsStr, &priceStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrBranchServiceNotFound
		}
		return nil, nil, fmt.Errorf("query failed: %w", err)
	}

	if !detailsStr.Valid || !priceStr.Valid {
		return nil, nil, ErrDetailsNotFound
	}

	var detailsMap map[string]int
	if detailsStr.String == "" || detailsStr.String == "null" {
		detailsMap = make(map[string]int)
	} else {
		if err := json.Unmarshal([]byte(detailsStr.String), &detailsMap); err != nil {
			return nil, nil, fmt.Errorf("failed to parse service details: %w", err)
		}
	}

	var priceMap map[string]float32
	if priceStr.String == "" || priceStr.String == "null" {
		priceMap = make(map[string]float32)
	} else {
		if err := json.Unmarshal([]byte(priceStr.String), &priceMap); err != nil {
			return nil, nil, fmt.Errorf("failed to parse service price: %w", err)
		}
	}

	detailsStruct := make([]*ServiceDetails, 0, len(detailsMap))
	for detail, duration := range detailsMap {
		detailsStruct = append(detailsStruct, &ServiceDetails{
			Detail:   detail,
			Duration: duration,
		})
	}

	priceStruct := make([]*ServPrice, 0, len(priceMap))
	for detail, price := range priceMap {
		priceStruct = append(priceStruct, &ServPrice{
			Detail: detail,
			Price:  price,
		})
	}

	return detailsStruct, priceStruct, nil
}

// Получение филала в определённом городе с определённой услугой
func (s *PostgresBranchStorage) GetBranchByCityServ(city string, serviceID string) ([]*BrancByCityServ, error) {
	rows, err := s.DB.Query(`
	SELECT 
	bs.id AS branch_service_id,
    b.id, 
    b.address,
    cmp.org_short_name
	FROM branches b
	JOIN branch_services bs ON b.id = bs.branch
	JOIN companies cmp ON b.inn_company = cmp.inn
	WHERE b.city = $1 AND bs.service = $2
`, city, serviceID)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var result []*BrancByCityServ
	for rows.Next() {
		var bcs BrancByCityServ
		var address, orgShortName sql.NullString

		// Сканирование UUID напрямую в поля типа uuid.UUID (поддерживается драйвером lib/pq)
		if err := rows.Scan(
			&bcs.BranchServId,
			&bcs.BranchId,
			&address,
			&orgShortName); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Обработка потенциально NULL-значений
		if address.Valid {
			bcs.Address = address.String
		}
		if orgShortName.Valid {
			bcs.CompanyName = orgShortName.String
		}

		result = append(result, &bcs)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return result, nil
}

// Создание филиала
func (s *PostgresBranchStorage) CreateBranch(branch Branch) (*Branch, error) {
	var id uuid.UUID
	err := s.DB.QueryRow(`
        INSERT INTO branches (city, address, inn_company, open_time, close_time)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, branch.City, branch.Address, branch.Inn, branch.OpenTime, branch.CloseTime).Scan(&id)
	log.Printf("%v", branch.OpenTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	created := &Branch{
		ID:        id,
		City:      branch.City,
		Address:   branch.Address,
		Inn:       branch.Inn,
		OpenTime:  branch.OpenTime,
		CloseTime: branch.CloseTime,
	}
	return created, nil
}

// Create добавляет новую запись в таблицу branch_services.
// Все поля, кроме BusyTime, обязательны. ID генерируется базой.
// Возвращает ошибку ErrBranchServiceExists, если пара branch+service уже существует.
// TODO(ylly): Вынести обрпботку ошибки  ErrBranchServiceExists в servise - через get
func (s *PostgresBranchStorage) CreateBranchServ(bs BranchServ) (*BranchServ, error) {
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

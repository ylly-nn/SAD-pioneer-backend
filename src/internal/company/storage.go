package company

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"src/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// ошибки которые возвращает company/storage
var (
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
	ErrBranchNotFound       = errors.New("branch not found")
	ErrBranchesNotFound     = errors.New("company has no branches")
	ErrBranchServNotFound   = errors.New("service by branch not found")
)

// CompanyStorage определяет методы для работы с компаниями
type CompanyStorage interface {
	GetPartUserByEmail(email string) (PartnersUsers, error)

	GetAll() ([]Company, error)

	GetCompanyByInn(inn string) (*Company, error)

	Create(Company) (*Company, error)

	Delete(inn string) error

	GetBranchesByInn(inn string) ([]*CompanyBranch, error)

	GetBranchByID(branchID uuid.UUID) (CompanyBranch, error)

	GetServicesByBranch(branchID uuid.UUID) ([]*ServiceInBranch, error)

	GetBranchServByID(branchServID uuid.UUID) (BranchServ, error)
}

// PostgresCompanyStorage реализует CompanyStorage для PostgreSQL.
type PostgresCompanyStorage struct {
	*db.Storage
}

// NewPostgresCompanyStorage создаёт новый экземпляр PostgresCompanyStorage.
func NewPostgresCompanyStorage(sqlDB *sql.DB) *PostgresCompanyStorage {
	return &PostgresCompanyStorage{Storage: db.NewStorage(sqlDB)}
}

// Получение из бд branch_serv по id если нет, ошибка - ErrBranchServNotFound
func (s *PostgresCompanyStorage) GetBranchServByID(branchServID uuid.UUID) (BranchServ, error) {
	var bs BranchServ
	var details []byte

	row := s.DB.QueryRow(`
        SELECT id, branch, service, service_detalis
        FROM branch_services
        WHERE id = $1
    `, branchServID)

	err := row.Scan(&bs.ID, &bs.Branch, &bs.Service, &details)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BranchServ{}, ErrBranchServNotFound
		}
		return BranchServ{}, err
	}

	bs.ServiceDetails = json.RawMessage(details)
	return bs, nil
}

// Получение из бд филиала по его id
// Если не найдено ErrBranchNotFound
func (s *PostgresCompanyStorage) GetBranchByID(branchID uuid.UUID) (CompanyBranch, error) {
	var branch CompanyBranch

	row := s.DB.QueryRow(`
        SELECT id, city, address, inn_company, open_time, close_time
        FROM branches
        WHERE id = $1
    `, branchID)

	err := row.Scan(&branch.ID, &branch.City, &branch.Address, &branch.Inn, &branch.OpenTime, &branch.CloseTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CompanyBranch{}, ErrBranchNotFound
		}
		return CompanyBranch{}, fmt.Errorf("scan branch: %w", err)
	}

	return branch, nil
}

// Получение сиска услуг определённого филиала
func (s *PostgresCompanyStorage) GetServicesByBranch(branchID uuid.UUID) ([]*ServiceInBranch, error) {
	rows, err := s.DB.Query(`
        SELECT 
            bs.id AS branch_service_id, 
            s.id AS service_id, 
            s.name AS service_name
        FROM branch_services bs
        JOIN services s ON s.id = bs.service
        WHERE bs.branch = $1
    `, branchID)
	if err != nil {
		return nil, fmt.Errorf("query services by branch %v: %w", branchID, err)
	}
	defer rows.Close()

	var services []*ServiceInBranch
	for rows.Next() {
		var serv ServiceInBranch
		if err := rows.Scan(&serv.BranchServId, &serv.ServiceId, &serv.ServiceName); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}
		services = append(services, &serv)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return services, nil
}

// GetBrancesByCompany получение филиалов по инн компании
// возвращает масиив с филиалами
func (s *PostgresCompanyStorage) GetBranchesByInn(inn string) ([]*CompanyBranch, error) {
	rows, err := s.DB.Query(`
        SELECT id, city, address, inn_company, open_time, close_time
        FROM branches
        WHERE inn_company = $1
    `, inn)
	if err != nil {
		return nil, fmt.Errorf("query branches by company inn %s: %w", inn, err)
	}
	defer rows.Close()

	var branches []*CompanyBranch
	for rows.Next() {
		var b CompanyBranch
		// Сканируем напрямую – TimeOnly сам разберёт строку
		err := rows.Scan(&b.ID, &b.City, &b.Address, &b.Inn, &b.OpenTime, &b.CloseTime)
		if err != nil {
			return nil, fmt.Errorf("scan branch: %w", err)
		}
		branches = append(branches, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	if len(branches) == 0 {
		return nil, ErrBranchesNotFound
	}
	return branches, nil
}

// GetPartUserByEmail - получение партнёра по email
func (s *PostgresCompanyStorage) GetPartUserByEmail(email string) (PartnersUsers, error) {
	var psrtnerUser PartnersUsers
	row := s.DB.QueryRow(`
        SELECT email, inn
        FROM partners_users
        WHERE email = $1
    `, email)

	err := row.Scan(&psrtnerUser.Email, &psrtnerUser.Inn)
	if err != nil {
		// Если запись не найдена, возвращаем пустую структуру без ошибки.
		if errors.Is(err, sql.ErrNoRows) {
			return PartnersUsers{}, nil
		}
		// Любая другая ошибка возвращается как есть.
		return PartnersUsers{}, fmt.Errorf("failed to scan partner user: %w", err)
	}

	return psrtnerUser, nil

}

// GetAll возвращает список всех компаний
func (s *PostgresCompanyStorage) GetAll() ([]Company, error) {
	rows, err := s.DB.Query(`SELECT inn, kpp, ogrn, org_name, org_short_name FROM companies`)
	if err != nil {
		return nil, fmt.Errorf("query all companies: %w", err)
	}
	defer rows.Close()

	var company []Company
	for rows.Next() {
		var cmp Company
		if err := rows.Scan(&cmp.INN, &cmp.KPP, &cmp.OGRN, &cmp.OrgName, &cmp.OrgShortName); err != nil {
			return nil, fmt.Errorf("scan company: %w", err)
		}
		company = append(company, cmp)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return company, nil
}

// // Delete удаляет компанию по inn
// // Если компания не найдена, возвращает companyNotFound
func (s *PostgresCompanyStorage) Delete(inn string) error {
	result, err := s.DB.Exec(`DELETE FROM companies WHERE inn = $1`, inn)
	if err != nil {
		return fmt.Errorf("delete company %v: %w", inn, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for delete %v: %w", inn, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete company %v: %w", inn, ErrCompanyNotFound)
	}
	return nil
}

// GetByInn возвращает компанию по её ИНН.
// Если компания не найдена, возвращает ошибку ErrCompanyNotFound.
func (s *PostgresCompanyStorage) GetCompanyByInn(inn string) (*Company, error) {
	var cmp Company
	var kpp, ogrn, orgName, orgShortName sql.NullString

	err := s.DB.QueryRow(`
		SELECT inn, kpp, ogrn, org_name, org_short_name
		FROM companies
		WHERE inn = $1
	`, inn).Scan(
		&cmp.INN,
		&kpp,
		&ogrn,
		&orgName,
		&orgShortName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("query company by inn %s: %w", inn, err)
	}

	if kpp.Valid {
		cmp.KPP = &kpp.String
	}
	if ogrn.Valid {
		cmp.OGRN = &ogrn.String
	}
	if orgName.Valid {
		cmp.OrgName = &orgName.String
	}
	if orgShortName.Valid {
		cmp.OrgShortName = &orgShortName.String
	}

	return &cmp, nil
}

// Create добавляет новую компанию в базу данных.
// Если компания с таким ИНН уже существует, возвращает ErrCompanyAlreadyExists.
func (s *PostgresCompanyStorage) Create(company Company) (*Company, error) {
	_, err := s.DB.Exec(`
        INSERT INTO companies (inn, kpp, ogrn, org_name, org_short_name)
        VALUES ($1, $2, $3, $4, $5)
    `, company.INN, company.KPP, company.OGRN, company.OrgName, company.OrgShortName)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, fmt.Errorf("create company: %w", ErrCompanyAlreadyExists)
		}
		return nil, fmt.Errorf("create company: %w", err)
	}

	return &company, nil
}

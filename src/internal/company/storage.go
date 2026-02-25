package company

import (
	"database/sql"
	"errors"
	"fmt"

	"src/internal/db"

	"github.com/jackc/pgx/v5/pgconn"
)

// ошибки которые возвращает company/storage
var (
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
)

// CompanyStorage определяет методы для работы с компаниями
type CompanyStorage interface {
	GetAll() ([]Company, error)

	GetCompanyByInn(inn string) (*Company, error)

	Create(Company) (*Company, error)

	Delete(inn string) error
}

// PostgresCompanyStorage реализует CompanyStorage для PostgreSQL.
type PostgresCompanyStorage struct {
	*db.Storage
}

// NewPostgresCompanyStorage создаёт новый экземпляр PostgresCompanyStorage.
func NewPostgresCompanyStorage(sqlDB *sql.DB) *PostgresCompanyStorage {
	return &PostgresCompanyStorage{Storage: db.NewStorage(sqlDB)}
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

package partners

import (
	"database/sql"
	"errors"
	"fmt"
	"src/internal/auth"
	"time"
)

// Интерфейс для работы с пользователями
type UserStorage interface {
	GetByEmail(email string) (*auth.User, error)
	Create(user *auth.User) error
	UpdatePassword(email, password string) error
}

// CompanyStorage интерфейс для работы с компаниями
type CompanyStorage interface {
	GetByINN(inn string) (*Company, error)
	Exists(inn string) (bool, error)
}

// PartnerRequestStorage интерфейс для работы с заявками
type PartnerRequestStorage interface {
	Create(req *PartnerRequest) error
	GetByINN(inn string) (*PartnerRequest, error)
	GetByUserEmail(email string) (*PartnerRequest, error)
	GetPartUserByEmail(email string) (PartnersUsers, error)
	Delete(inn string) error
}

// PartnersUsersStorage интерфейс для работы с таблицей partners_users
type PartnersUsersStorage interface {
	Create(email, inn string) error
}

type PostgresPartnerRequestStorage struct {
	db *sql.DB
}

func NewPostgresPartnerRequestStorage(db *sql.DB) *PostgresPartnerRequestStorage {
	return &PostgresPartnerRequestStorage{db: db}
}

// Создание заявки на подключение организации
func (s *PostgresPartnerRequestStorage) Create(req *PartnerRequest) error {
	query := `
        INSERT INTO part_req (
            status, user_email, inn, kpp, ogrn, org_name, org_short_name,
            name, surname, patronymic, email, phone_number, info, created_at, last_used
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `

	_, err := s.db.Exec(
		query,
		"new",
		req.UserEmail, req.INN, req.KPP, req.OGRN, req.OrgName, req.OrgShortName,
		req.Name, req.Surname, req.Patronymic, req.Email, req.Phone, req.Info, time.Now(), time.Now(),
	)

	return err
}

// Получение информации из заявки по ИНН
func (s *PostgresPartnerRequestStorage) GetByINN(inn string) (*PartnerRequest, error) {
	var req PartnerRequest
	query := `SELECT status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                     name, surname, patronymic, email, phone_number, info, created_at, last_used 
              FROM part_req WHERE inn = $1`

	err := s.db.QueryRow(query, inn).Scan(
		&req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
		&req.OrgName, &req.OrgShortName,
		&req.Name, &req.Surname, &req.Patronymic,
		&req.Email, &req.Phone, &req.Info, &req.CreatedAt, &req.LastUsed,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// Удаление заявки по ИНН (для отката)
func (s *PostgresPartnerRequestStorage) Delete(inn string) error {
	query := `DELETE FROM part_req WHERE inn = $1`
	result, err := s.db.Exec(query, inn)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("partner request with inn %s not found", inn)
	}
	return nil
}

type PostgresCompanyStorage struct {
	db *sql.DB
}

func NewPostgresCompanyStorage(db *sql.DB) *PostgresCompanyStorage {
	return &PostgresCompanyStorage{db: db}
}

// Получение информации о компании по ИНН
func (s *PostgresCompanyStorage) GetByINN(inn string) (*Company, error) {
	var company Company
	query := `SELECT inn, kpp, ogrn, org_name, org_short_name FROM companies WHERE inn = $1`

	err := s.db.QueryRow(query, inn).Scan(
		&company.INN, &company.KPP, &company.OGRN,
		&company.OrgName, &company.OrgShortName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// Проверка на существование компании
func (s *PostgresCompanyStorage) Exists(inn string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE inn = $1)`

	err := s.db.QueryRow(query, inn).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetByUserEmail - получение заявки по email пользователя
func (s *PostgresPartnerRequestStorage) GetByUserEmail(email string) (*PartnerRequest, error) {
	var req PartnerRequest
	query := `SELECT status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                     name, surname, patronymic, email, phone_number, info, created_at, last_used  
              FROM part_req WHERE user_email = $1`

	err := s.db.QueryRow(query, email).Scan(
		&req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
		&req.OrgName, &req.OrgShortName,
		&req.Name, &req.Surname, &req.Patronymic,
		&req.Email, &req.Phone, &req.Info, &req.CreatedAt, &req.LastUsed,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// GetPartUserByEmail - получение партнёра по email
func (s *PostgresPartnerRequestStorage) GetPartUserByEmail(user_email string) (PartnersUsers, error) {
	var partnerUser PartnersUsers

	query := `SELECT inn, user_email FROM part_req WHERE user_email = $1`

	err := s.db.QueryRow(query, user_email).Scan(&partnerUser.Email, &partnerUser.Inn)

	if err != nil {
		// Если запись не найдена, возвращаем пустую структуру без ошибки.
		if errors.Is(err, sql.ErrNoRows) {
			return PartnersUsers{}, nil
		}
		// Любая другая ошибка возвращается как есть.
		return PartnersUsers{}, fmt.Errorf("failed to scan partner user: %w", err)
	}

	return partnerUser, nil

}

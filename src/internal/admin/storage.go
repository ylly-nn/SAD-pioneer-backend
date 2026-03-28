package admin

import (
	"database/sql"
	"errors"
	"fmt"

	"src/internal/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// Интерфейс для работы с пользователями
type UserStorage interface {
	GetByEmail(email string) (*auth.User, error)
	Create(user *auth.User) error
	UpdatePassword(email, password string) error
}

// PartnerRequestStorage интерфейс для работы с заявками
type PartnerRequestStorage interface {
	Create(req *PartnerRequest) error
	GetByID(id uuid.UUID) (*PartnerRequest, error)
	GetByStatus(status string) ([]*PartnerRequest, error)
	GetPending() ([]*PartnerRequest, error)
	GetAll() ([]*PartnerRequest, error)
	UpdateStatus(id uuid.UUID, status string) error
	//Delete(inn string) error
}

// CompanyStorage интерфейс для работы с компаниями
type CompanyStorage interface {
	Create(company *Company) error
	GetByINN(inn string) (*Company, error)
	Exists(inn string) (bool, error)
}

// PartnersUsersStorage интерфейс для работы с таблицей partners_users
type PartnersUsersStorage interface {
	Create(email, inn string) error
}

// AdminStorage интерфейс для работы с админами
type AdminStorage interface {
	IsAdmin(email string) (bool, error)
	GetByEmail(email string) (*Admin, error)
	Create(email, name, surname string) error
}

// Admin структура для таблицы admin
type Admin struct {
	Email   string `json:"email" db:"email"`
	Name    string `json:"name" db:"name"`
	Surname string `json:"surname" db:"surname"`
}

// Create создаёт нового администратора
func (s *PostgresAdminStorage) Create(email, name, surname string) error {
	query := `INSERT INTO admin (email, name, surname) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(query, email, name, surname)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("admin with email %s already exists", email)
		}
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

// PostgresAdminStorage реализация для PostgreSQL
type PostgresAdminStorage struct {
	db *sql.DB
}

func NewPostgresAdminStorage(db *sql.DB) *PostgresAdminStorage {
	return &PostgresAdminStorage{db: db}
}

// Проверка на админа
func (s *PostgresAdminStorage) IsAdmin(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM admin WHERE email = $1)`

	err := s.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return exists, nil
}

// Получение информации об админе по email
func (s *PostgresAdminStorage) GetByEmail(email string) (*Admin, error) {
	var admin Admin
	query := `SELECT email, name, surname FROM admin WHERE email = $1`

	err := s.db.QueryRow(query, email).Scan(&admin.Email, &admin.Name, &admin.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
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
            name, surname, patronymic, email, phone_number, info
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err := s.db.Exec(
		query,
		"new",
		req.UserEmail, req.INN, req.KPP, req.OGRN, req.OrgName, req.OrgShortName,
		req.Name, req.Surname, req.Patronymic, req.Email, req.Phone, req.Info,
	)

	return err
}

// Получение информации из заявки по ID
func (s *PostgresPartnerRequestStorage) GetByID(id uuid.UUID) (*PartnerRequest, error) {
	var req PartnerRequest
	query := `SELECT id, status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                     name, surname, patronymic, email, phone_number, info 
              FROM part_req WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&req.ID, &req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
		&req.OrgName, &req.OrgShortName,
		&req.Name, &req.Surname, &req.Patronymic,
		&req.Email, &req.Phone, &req.Info,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// Обновление статуса у заявки
func (s *PostgresPartnerRequestStorage) UpdateStatus(id uuid.UUID, status string) error {
	query := `UPDATE part_req SET status = $1 WHERE id = $2`
	result, err := s.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("partner request with id %s not found", id)
	}
	return nil
}

// Получение заявок в работе
func (s *PostgresPartnerRequestStorage) GetPending() ([]*PartnerRequest, error) {
	rows, err := s.db.Query(`SELECT id, status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                                     name, surname, patronymic, email, phone_number, info 
                              FROM part_req WHERE status = 'pending'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*PartnerRequest
	for rows.Next() {
		var req PartnerRequest
		err := rows.Scan(
			&req.ID, &req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
			&req.OrgName, &req.OrgShortName,
			&req.Name, &req.Surname, &req.Patronymic,
			&req.Email, &req.Phone, &req.Info,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

// Получение всех заявок
func (s *PostgresPartnerRequestStorage) GetAll() ([]*PartnerRequest, error) {
	rows, err := s.db.Query(`SELECT id, status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                                     name, surname, patronymic, email, phone_number, info 
                              FROM part_req`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*PartnerRequest
	for rows.Next() {
		var req PartnerRequest
		err := rows.Scan(
			&req.ID, &req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
			&req.OrgName, &req.OrgShortName,
			&req.Name, &req.Surname, &req.Patronymic,
			&req.Email, &req.Phone, &req.Info,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

// Получение заявок с определенным статусом
func (s *PostgresPartnerRequestStorage) GetByStatus(status string) ([]*PartnerRequest, error) {
	rows, err := s.db.Query(`SELECT id, status, user_email, inn, kpp, ogrn, org_name, org_short_name,
                                     name, surname, patronymic, email, phone_number, info 
                              FROM part_req WHERE status = $1`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*PartnerRequest
	for rows.Next() {
		var req PartnerRequest
		err := rows.Scan(
			&req.ID, &req.Status, &req.UserEmail, &req.INN, &req.KPP, &req.OGRN,
			&req.OrgName, &req.OrgShortName,
			&req.Name, &req.Surname, &req.Patronymic,
			&req.Email, &req.Phone, &req.Info,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
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

// Создание компании
func (s *PostgresCompanyStorage) Create(company *Company) error {
	query := `
		INSERT INTO companies (inn, kpp, ogrn, org_name, org_short_name)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(
		query,
		company.INN, company.KPP, company.OGRN,
		company.OrgName, company.OrgShortName,
	)
	return err
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

// PostgresPartnersUsersStorage реализация для PostgreSQL
type PostgresPartnersUsersStorage struct {
	db *sql.DB
}

func NewPostgresPartnersUsersStorage(db *sql.DB) *PostgresPartnersUsersStorage {
	return &PostgresPartnersUsersStorage{db: db}
}

// Создание нулевого пользователя от организации после одобрения заявки
func (s *PostgresPartnersUsersStorage) Create(email, inn string) error {
	query := `INSERT INTO partners_users (email, inn) VALUES ($1, $2)`

	_, err := s.db.Exec(query, email, inn)
	if err != nil {
		return fmt.Errorf("failed to insert into partners_users: %w", err)
	}

	return nil
}

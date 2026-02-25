package company

import (
	"errors"
	"fmt"
)

type CompanyManager struct {
	storage CompanyStorage
}

func NewCompanyManager(storage CompanyStorage) *CompanyManager {
	return &CompanyManager{storage: storage}
}

// GetAllCompanies - возвращает список всех компаний
func (m *CompanyManager) GetAllCompanies() ([]Company, error) {
	companies, err := m.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all companies: %w", err)
	}
	return companies, nil
}

// DeleteCompany удаляет компанию  по инн
// Если компания не найдена, возвращает ошибку ErrCompanyNotFound.
func (m *CompanyManager) DeleteCompany(inn string) error {
	return m.storage.Delete(inn)
}

// GetCompanyByInn возвращает компанию  по инн
// Если компания не найдена, возвращает ошибку ErrCompanyNotFound.
func (m *CompanyManager) GetCompanyByInn(inn string) (*Company, error) {
	company, err := m.storage.GetCompanyByInn(inn)
	if err != nil {
		return nil, fmt.Errorf("get company by inn: %w", err)
	}
	return company, nil
}

// CreateCompany создаёт новую компанию.
// Все поля Company должны быть заполнены (не nil) и соответствовать форматам.
func (m *CompanyManager) CreateCompany(company Company) (*Company, error) {
	if company.INN == "" {
		return nil, errors.New("INN is required")
	}
	if len(company.INN) != 10 && len(company.INN) != 12 {
		return nil, errors.New("INN must be 10 or 12 characters")
	}

	if company.KPP == nil || *company.KPP == "" {
		return nil, errors.New("KPP is required and cannot be empty")
	}
	if len(*company.KPP) != 9 {
		return nil, errors.New("KPP must be 9 characters")
	}

	if company.OGRN == nil || *company.OGRN == "" {
		return nil, errors.New("OGRN is required and cannot be empty")
	}
	if len(*company.OGRN) != 13 {
		return nil, errors.New("OGRN must be 13 characters")
	}

	if company.OrgName == nil || *company.OrgName == "" {
		return nil, errors.New("organization name is required and cannot be empty")
	}
	if company.OrgShortName == nil || *company.OrgShortName == "" {
		return nil, errors.New("short organization name is required and cannot be empty")
	}

	created, err := m.storage.Create(company)
	if err != nil {
		return nil, fmt.Errorf("create company: %w", err)
	}
	return created, nil
}

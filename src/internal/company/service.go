package company

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

var (
	ErrUserNotPartner         = errors.New("the user does not have a company")
	ErrBranchNotInCompany     = errors.New("no access to the company that owns the branch")
	ErrBranchServNotAvailable = errors.New("service in the branch not available to the user")
	ErrBranchServIsNull       = errors.New("details for service in the branch not found")
	ErrServiceDetailsInvalid  = errors.New("invalid service details format")
)

// CompanyManager содержит бизнес-логику для работы с компаниями.
type CompanyManager struct {
	storage     CompanyStorage
	userStorage UserStorage
}

// NewCompanyManager создаёт новый экземпляр CompanyManager.
func NewCompanyManager(storage CompanyStorage, userStorage UserStorage) *CompanyManager {
	return &CompanyManager{storage: storage,
		userStorage: userStorage}
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

// GetCompany возвращает компанию  по инн
// Если компания не найдена, возвращает ошибку ErrCompanyNotFound.
func (m *CompanyManager) GetCompany(email string) (*Company, error) {
	userIsPartner, err := m.UserIsPartner(email)

	if userIsPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	company, err := m.storage.GetCompanyByInn(userIsPartner.Inn)
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

// GetBranchesByEmail - получение филлиалов компании инн которой получаем по email
func (m *CompanyManager) GetBranchesByEmail(email string) ([]*CompanyBranch, error) {
	userIsPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, err
	}
	if userIsPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	branches, err := m.storage.GetBranchesByInn(userIsPartner.Inn)

	if err != nil {
		return nil, err
	}
	return branches, nil
}

// получение филиала компании по id и email пользователя принадлежащего компании
func (m *CompanyManager) GetBranchByIdEmail(branch_id uuid.UUID, email string) (CompanyBranchWithServ, error) {

	isPartner, err := m.UserIsPartner(email)
	if isPartner.IsPartner != true {
		return CompanyBranchWithServ{}, ErrUserNotPartner
	}

	branch, err := m.storage.GetBranchByID(branch_id)
	if errors.Is(err, ErrBranchNotFound) {
		return CompanyBranchWithServ{}, ErrBranchNotInCompany
	}

	if err != nil {
		return CompanyBranchWithServ{}, err
	}

	if branch.Inn != isPartner.Inn {
		return CompanyBranchWithServ{}, ErrBranchNotInCompany
	}

	serv, err := m.storage.GetServicesByBranch(branch_id)
	if err != nil {
		return CompanyBranchWithServ{}, err
	}

	var bws CompanyBranchWithServ

	bws.City = branch.City
	bws.Address = branch.Address
	bws.OpenTime = branch.OpenTime
	bws.OpenTime = branch.CloseTime
	bws.Services = serv

	return bws, nil
}

// Получение деталей ууслуги определёного филиала
// Возвращаемые ошибки:  ErrUserNotPartner, ErrBranchesNotFound
// ErrBranchServNotFound, ErrBranchServNotAvailable
// ErrBranchServNotAvailable, ErrServiceDetailsInvalid
func (m *CompanyManager) GetServDetailsByBranchServId(branchServID uuid.UUID, email string) ([]*CompanyServDetailsResponse, error) {
	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return []*CompanyServDetailsResponse{}, err
	}
	if isPartner.IsPartner != true {
		return []*CompanyServDetailsResponse{}, ErrUserNotPartner
	}

	branchServ, err := m.storage.GetBranchServByID(branchServID)
	if err != nil {
		return []*CompanyServDetailsResponse{}, err
	}

	companyBranch, err := m.storage.GetBranchesByInn(isPartner.Inn)
	if err != nil {
		return []*CompanyServDetailsResponse{}, err
	}

	//создание массива из id филиалов компании
	var branchIDs []uuid.UUID
	for _, branch := range companyBranch {
		branchIDs = append(branchIDs, branch.ID)
	}

	//проверка есть ли в списке id филиалов id филиала в branchServ
	found := slices.Contains(branchIDs, branchServ.Branch)

	if !found {
		return []*CompanyServDetailsResponse{}, ErrBranchServNotAvailable
	}

	if len(branchServ.ServiceDetails) == 0 || string(branchServ.ServiceDetails) == "null" {
		return []*CompanyServDetailsResponse{}, ErrBranchServIsNull
	}

	var detailsMap map[string]int
	if err := json.Unmarshal(branchServ.ServiceDetails, &detailsMap); err != nil {
		return nil, ErrServiceDetailsInvalid
	}

	var result []*CompanyServDetailsResponse
	for detail, duration := range detailsMap {
		result = append(result, &CompanyServDetailsResponse{
			Detail:   detail,
			Duration: duration,
		})
	}
	return result, nil

}

// Проверка что у пользователя есть организация
// Если у пользователя организация есть вернётся IsParnersUsers{IsPartner: true, Inn: "строка с инн"
// Если у пользователя организация нет вернётся IsParnersUsers{IsPartner: false, Inn: ничего
func (m *CompanyManager) UserIsPartner(email string) (IsPartnersUsers, error) {
	partUser, err := m.storage.GetPartUserByEmail(email)
	if err != nil {
		return IsPartnersUsers{IsPartner: false}, fmt.Errorf("failed to check if user is partner: %w", err)
	}

	if partUser.Email == "" {
		return IsPartnersUsers{IsPartner: false}, nil
	}

	return IsPartnersUsers{
		IsPartner: true,
		Inn:       partUser.Inn,
	}, nil
}

// AddUserToCompany добавляет нового пользователя в компанию
func (m *CompanyManager) AddUserToCompany(userEmail, newUserEmail string) error {
	// Проверка, что добавляющий пользователь есть в компании
	userIsPartner, err := m.UserIsPartner(userEmail)
	if err != nil {
		return fmt.Errorf("failed to check user email: %w", err)
	}
	if !userIsPartner.IsPartner {
		return ErrUserNotPartner
	}

	inn := userIsPartner.Inn

	// Проверка, что компания с ИНН существует и пользователь находится в компании
	if userIsPartner.Inn != inn {
		return errors.New("user does not have access to this company")
	}

	company, err := m.storage.GetCompanyByInn(inn)
	if err != nil {
		return ErrCompanyNotFound
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Проверка, существует ли новый пользователь в системе
	user, err := m.userStorage.GetByEmail(newUserEmail)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Проверка, есть ли пользователь в компании
	existingPartner, err := m.storage.GetPartUserByEmail(newUserEmail)
	if err != nil {
		return fmt.Errorf("failed to check if user is already partner: %w", err)
	}
	if existingPartner.Email != "" {
		return errors.New("user is already a partner")
	}

	if err := m.storage.AddUserToPartners(newUserEmail, inn); err != nil {
		return fmt.Errorf("failed to add user to partners: %w", err)
	}

	return nil
}

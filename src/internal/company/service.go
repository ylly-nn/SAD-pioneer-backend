package company

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"src/internal/city"
	"src/internal/timeparsing"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrUserNotPartner                = errors.New("the user does not have a company")
	ErrBranchNotInCompany            = errors.New("no access to the company that owns the branch")
	ErrBranchServNotAvailable        = errors.New("service in the branch not available to the user")
	ErrBranchServIsNull              = errors.New("details for service in the branch not found")
	ErrServiceDetailsInvalid         = errors.New("invalid service details format")
	ErrEmptyCity                     = errors.New("city cannot be empty")
	ErrInvalidCity                   = errors.New("city is not in the list of Russian cities")
	ErrOrderNotAvailable             = errors.New("order not available to the user")
	ErrUpdateStatus                  = errors.New("status cannstatus cannot be changed to the selected oneot be changed to the selected one")
	ErrStatus                        = errors.New("status must be approve or reject")
	ErrBranchServDetailAlreadyExists = errors.New("detail for this service in the branch already exists")
	ErrInvalidDuration               = errors.New("invalid duration")
	ErrDetailNotFound                = errors.New("detail in branch service not found")
)

var hyphenSpaces = regexp.MustCompile(`\s*-\s*`)

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
func (m *CompanyManager) GetServDetailsByBranchServId(branchServID uuid.UUID, email string) ([]*ServUpdateResponse, error) {
	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, err
	}
	if isPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	branchServ, err := m.storage.GetBranchServByID(branchServID)
	if err != nil {
		return nil, err
	}

	companyBranch, err := m.storage.GetBranchesByInn(isPartner.Inn)
	if err != nil {
		return nil, err
	}

	//создание массива из id филиалов компании
	var branchIDs []uuid.UUID
	for _, branch := range companyBranch {
		branchIDs = append(branchIDs, branch.ID)
	}

	//проверка есть ли в списке id филиалов id филиала в branchServ
	found := slices.Contains(branchIDs, branchServ.Branch)

	if !found {
		return nil, ErrBranchServNotAvailable
	}

	if len(branchServ.ServiceDetails) == 0 {
		return nil, ErrBranchServIsNull
	}

	result := make([]*ServUpdateResponse, len(branchServ.ServiceDetails))
	for i := range branchServ.ServiceDetails {
		result[i] = &branchServ.ServiceDetails[i]
	}

	return result, nil

}

// GetCompanyOrders возвращает заказы по всем филиалам компании партнёра
func (m *CompanyManager) GetCompanyOrders(email string) ([]*CompanyBranchOrderResponse, error) {
	// Проверяем, является ли пользователь партнёром
	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, fmt.Errorf("check partner status: %w", err)
	}
	if !isPartner.IsPartner {
		return nil, ErrUserNotPartner
	}

	// Получаем все филиалы компании по ИНН
	branches, err := m.storage.GetBranchesByInn(isPartner.Inn)
	if err != nil {
		if errors.Is(err, ErrBranchesNotFound) {
			return []*CompanyBranchOrderResponse{}, ErrBranchNotFound
		}
		return nil, fmt.Errorf("get branches: %w", err)
	}

	var result []*CompanyBranchOrderResponse

	// Для каждого филиала получаем заказы
	for _, branch := range branches {
		orders, err := m.storage.GetOrdersByBranch(branch.ID)
		if err != nil {
			// Если заказов нет, пропускаем филиал
			if errors.Is(err, ErrOrderNotFound) {
				continue
			}
			// Другая ошибка — прерываем выполнение
			return nil, fmt.Errorf("get orders for branch %s: %w", branch.ID, err)
		}

		// Если заказы есть, добавляем филиал в ответ
		result = append(result, &CompanyBranchOrderResponse{
			BranchID: branch.ID,
			City:     branch.City,
			Address:  branch.Address,
			Orders:   orders,
		})
	}
	if len(result) == 0 {
		return nil, ErrOrderNotFound
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

// AddBranchToCompany добавляет новый филиал в компанию
func (m *CompanyManager) AddBranchToCompany(userEmail, cityName, address string, open_time, close_time timeparsing.TimeOnly) error {
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

	if cityName == "" {
		return ErrEmptyCity
	}

	fields := strings.Fields(cityName)
	if len(fields) == 0 {
		return ErrEmptyCity
	}
	cityName = strings.Join(fields, " ")

	cityName = hyphenSpaces.ReplaceAllString(cityName, "-")

	canonicalCity, ok := city.ValidCitiesMap[strings.ToLower(cityName)]
	if !ok {
		return ErrInvalidCity
	}

	// Проверка на уникальность адреса для этой компании
	exists, err := m.storage.CheckBranchAddressExists(inn, address, canonicalCity)
	if err != nil {
		return fmt.Errorf("failed to check branch address: %w", err)
	}
	if exists {
		return errors.New("branch with this address already exists for this company")
	}

	if err := m.storage.AddNewBranchToCompany(canonicalCity, address, inn, open_time, close_time); err != nil {
		return fmt.Errorf("failed to add branch to company: %w", err)
	}

	return nil
}

// AddServiceToBranch добавляет новую услугу в филиал
func (m *CompanyManager) AddServiceToBranch(userEmail string, branch_id, service_id uuid.UUID) error {
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
		if errors.Is(err, ErrCompanyNotFound) {
			return ErrCompanyNotFound
		}
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Проверка, что филиал просто есть
	branch, err := m.storage.GetBranchByID(branch_id)

	if err != nil {
		if errors.Is(err, ErrBranchNotFound) {
			return ErrBranchNotInCompany
		}
	}

	// Проверка, что филиал в компании
	if branch.Inn != inn {
		return ErrBranchNotInCompany
	}

	// Проверка, что услуга существует
	service, err := m.storage.GetServiceByID(service_id)
	if err != nil {
		// если услуги нет - 403
		return ErrBranchNotInCompany
	}
	if service == nil {
		return ErrBranchNotInCompany
	}

	// Проверка, что услуга уже не добавлена в филиал
	exists, err := m.storage.CheckServiceInBranchExists(branch_id, service_id)
	if err != nil {
		return fmt.Errorf("failed to check service in branch: %w", err)
	}
	if exists {
		return errors.New("service already exists in this branch")
	}

	// Добавление услуги в филиал
	if err := m.storage.AddServiceToBranch(branch_id, service_id); err != nil {
		return fmt.Errorf("failed to add service to branch: %w", err)
	}

	return nil
}

// обновляет статус заказа со стороны организации с проверками доступа
func (m *CompanyManager) UpdateOrderStatus(email string, orderId uuid.UUID, statusStr string) (*CompanyOrder, error) {
	var status OrderStatus

	if statusStr != "approve" && statusStr != "reject" {
		return nil, ErrStatus
	}

	if statusStr == "approve" {
		status = OrderStatusApprove
	}

	if statusStr == "reject" {
		status = OrderStatusReject
	}

	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, err
	}

	if isPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	companyOrders, err := m.GetCompanyOrders(email)

	if err != nil {
		return nil, err
	}

	// Поиск нужного заказа и сбор ID
	var targetOrder *CompanyOrder
	var allOrderIDs []uuid.UUID
	for _, branch := range companyOrders {
		for _, order := range branch.Orders {
			allOrderIDs = append(allOrderIDs, order.ID)
			if order.ID == orderId {
				targetOrder = order
			}
		}
	}

	log.Printf("Order IDs: %v", allOrderIDs)

	if targetOrder == nil {
		return nil, ErrOrderNotAvailable
	}

	// Проверка допустимости перехода статуса
	currentStatus := targetOrder.Status
	switch currentStatus {
	case OrderStatusCreate:
		if status != OrderStatusApprove && status != OrderStatusReject {
			return nil, ErrUpdateStatus
		}
	case OrderStatusApprove:
		if status != OrderStatusReject {
			return nil, ErrUpdateStatus
		}
	default:
		// reject или неизвестный статус - обновление запрещено
		return nil, ErrUpdateStatus
	}

	updatedOrder, err := m.storage.UpdateOrderStatus(orderId, status)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

// Добавляет деталь услуги из филиала по названию и длительности
func (m *CompanyManager) AddServiceDetail(branchServID uuid.UUID, email string, getDetail ServDetails, getPrices ServPrice) ([]*ServUpdateResponse, error) {

	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, err
	}

	if isPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	branchServ, err := m.storage.GetBranchServByID(branchServID)
	if err != nil {
		return nil, err
	}

	branch, err := m.storage.GetBranchByID(branchServ.Branch)
	if err != nil {
		return nil, err
	}

	if branch.Inn != isPartner.Inn {
		return nil, ErrBranchNotInCompany
	}

	dbDetails, dbPrice, err := m.storage.GetServiceDetailsAndPrice(branchServID)

	if err != nil {
		return nil, err
	}

	for _, existing := range dbDetails {
		if existing.Detail == getDetail.Detail {
			return nil, ErrBranchServDetailAlreadyExists
		}
	}

	dbDetails = append(dbDetails, &getDetail)
	dbPrice = append(dbPrice, &getPrices)

	detailsMap := make(map[string]int, len(dbDetails))
	for _, d := range dbDetails {
		detailsMap[d.Detail] = d.Duration
	}

	jsonRowDetails, err := json.Marshal(detailsMap)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[string]float64, len(dbPrice))
	for _, p := range dbPrice {
		// p.Price имеет тип float32, приводим к float64 для безопасного JSON
		priceMap[p.Detail] = float64(p.Price)
	}
	jsonRowPrices, err := json.Marshal(priceMap)
	if err != nil {
		return nil, err
	}

	if err := m.storage.UpdateServiceDetails(branchServID, jsonRowDetails, jsonRowPrices); err != nil {
		return nil, err
	}

	response := make([]*ServUpdateResponse, 0, len(dbDetails))
	for _, d := range dbDetails {
		priceVal, ok := priceMap[d.Detail]
		if !ok {
			return nil, fmt.Errorf("price not found for detail %s", d.Detail)
		}
		response = append(response, &ServUpdateResponse{
			Detail:   d.Detail,
			Duration: d.Duration,
			Price:    float32(priceVal),
		})
	}
	return response, nil

}

// DeleteServiceDetail удаляет деталь услуги из филиала по названию.
func (m *CompanyManager) DeleteServiceDetail(branchServID uuid.UUID, email string, nameDetail string) ([]*ServUpdateResponse, error) {

	isPartner, err := m.UserIsPartner(email)
	if err != nil {
		return nil, err
	}

	if isPartner.IsPartner != true {
		return nil, ErrUserNotPartner
	}

	branchServ, err := m.storage.GetBranchServByID(branchServID)
	if err != nil {
		return nil, err
	}

	branch, err := m.storage.GetBranchByID(branchServ.Branch)
	if err != nil {
		return nil, err
	}
	if branch.Inn != isPartner.Inn {
		return nil, ErrBranchNotInCompany
	}

	dbDetails, dbPrice, err := m.storage.GetServiceDetailsAndPrice(branchServID)
	if err != nil {
		return nil, err
	}

	// Ищем индекс удаляемой детали
	foundIndex := -1
	for i, d := range dbDetails {
		if d.Detail == nameDetail {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return nil, ErrDetailNotFound
	}

	// Удалить элемент из слайса (сохраняя порядок)
	dbDetails = append(dbDetails[:foundIndex], dbDetails[foundIndex+1:]...)

	priceIndex := -1
	for i, p := range dbPrice {
		if p.Detail == nameDetail {
			priceIndex = i
			break
		}
	}

	if priceIndex == -1 {
		return nil, fmt.Errorf("price not found for detail %s", nameDetail)
	}

	dbPrice = append(dbPrice[:priceIndex], dbPrice[priceIndex+1:]...)

	detailsMap := make(map[string]int, len(dbDetails))
	for _, d := range dbDetails {
		detailsMap[d.Detail] = d.Duration
	}

	//Сериализовать map в JSON
	jsonRowDetails, err := json.Marshal(detailsMap)
	if err != nil {
		return nil, fmt.Errorf("marshal service details: %w", err)
	}

	priceMap := make(map[string]float64, len(dbPrice))
	for _, p := range dbPrice {
		priceMap[p.Detail] = float64(p.Price)
	}
	jsonRowPrices, err := json.Marshal(priceMap)
	if err != nil {
		return nil, fmt.Errorf("marshal price: %w", err)
	}

	if err := m.storage.UpdateServiceDetails(branchServID, jsonRowDetails, jsonRowPrices); err != nil {
		return nil, err
	}

	response := make([]*ServUpdateResponse, 0, len(dbDetails))
	for _, d := range dbDetails {
		priceVal, ok := priceMap[d.Detail]
		if !ok {
			return nil, fmt.Errorf("price not found for detail %s", d.Detail)
		}
		response = append(response, &ServUpdateResponse{
			Detail:   d.Detail,
			Duration: d.Duration,
			Price:    float32(priceVal),
		})
	}

	return response, nil
}

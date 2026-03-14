package admin

import (
	"fmt"
	"time"

	configPkg "src/internal/config"
)

type Config struct {
	JWTSecretKey    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	VerificationTTL time.Duration
}

// Содержит бизнес-логику для работы с админскими маршрутами и заявками для организаций
type AdminManager struct {
	userStorage           UserStorage
	partnerRequestStorage PartnerRequestStorage
	companyStorage        CompanyStorage
	partnersUsersStorage  PartnersUsersStorage
	adminStorage          AdminStorage
	emailSender           configPkg.EmailSender
	config                Config
}

// Создает новый экземпляр сервиса
func NewAdminManager(
	userStorage UserStorage,
	partnerRequestStorage PartnerRequestStorage,
	companyStorage CompanyStorage,
	partnersUsersStorage PartnersUsersStorage,
	adminStorage AdminStorage,
	emailSender configPkg.EmailSender,
	config Config,
) *AdminManager {
	return &AdminManager{
		userStorage:           userStorage,
		partnerRequestStorage: partnerRequestStorage,
		companyStorage:        companyStorage,
		partnersUsersStorage:  partnersUsersStorage,
		adminStorage:          adminStorage,
		emailSender:           emailSender,
		config:                config,
	}
}

// Интерфейс для проверки админа
func (m *AdminManager) IsAdmin(email string) (bool, error) {
	return m.adminStorage.IsAdmin(email)
}

// Создание заявки партнера (доступно любому авторизованному пользователю)
func (s *AdminManager) CreatePartnerRequest(userEmail string, req *PartnerRequestRequest) error {
	// Проверка на существование пользователя
	user, err := s.userStorage.GetByEmail(userEmail)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Проверка на наличие заявки с таким же ИНН
	existing, _ := s.companyStorage.GetByINN(req.INN)
	if existing != nil {
		return fmt.Errorf("company with this INN already exists")
	}

	// Создание заявки
	partnerReq := &PartnerRequest{
		Status:       "new",
		UserEmail:    userEmail,
		INN:          req.INN,
		KPP:          req.KPP,
		OGRN:         req.OGRN,
		OrgName:      req.OrgName,
		OrgShortName: req.OrgShortName,
		Name:         req.Name,
		Surname:      req.Surname,
		Patronymic:   req.Patronymic,
		Email:        req.Email,
		Phone:        req.Phone,
		Info:         req.Info,
	}

	if err := s.partnerRequestStorage.Create(partnerReq); err != nil {
		return fmt.Errorf("failed to create partner request: %w", err)
	}

	return nil
}

// Смена статуса заявки с "новая" на "в работе" (new -> pending)
func (s *AdminManager) TakeRequestToWork(inn string) error {
	// Получение заявки по INN
	req, err := s.partnerRequestStorage.GetByINN(inn)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}
	if req == nil {
		return fmt.Errorf("request with inn %s not found", inn)
	}

	// Проверка, что заявка в статусе "new"
	if req.Status != "new" {
		return fmt.Errorf("request cannot be taken to work: current status is %s", req.Status)
	}

	// Обновление статус на "pending"
	if err := s.partnerRequestStorage.UpdateStatus(inn, "pending"); err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

// Одобрение заявки (только для админов)
func (s *AdminManager) ApprovePartnerRequest(inn string) error {
	// Получение заявки по INN
	req, err := s.partnerRequestStorage.GetByINN(inn)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}
	if req == nil {
		return fmt.Errorf("request with inn %s not found", inn)
	}
	if req.Status != "pending" {
		return fmt.Errorf("request already processed")
	}

	// Создание компании
	company := &Company{
		INN:          req.INN,
		KPP:          req.KPP,
		OGRN:         req.OGRN,
		OrgName:      req.OrgName,
		OrgShortName: req.OrgShortName,
	}

	if err := s.companyStorage.Create(company); err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	// Создание нулевого пользователя компании
	if err := s.partnersUsersStorage.Create(req.UserEmail, inn); err != nil {
		return fmt.Errorf("failed to create partner user record: %w", err)
	}

	// Обновление статуса заявки
	if err := s.partnerRequestStorage.UpdateStatus(inn, "approved"); err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

// Отклонение заявки
func (s *AdminManager) RejectPartnerRequest(inn string) error {
	// Получение заявки по INN
	req, err := s.partnerRequestStorage.GetByINN(inn)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}
	if req == nil {
		return fmt.Errorf("request with inn %s not found", inn)
	}

	// Проверка, что заявка в статусе "pending"
	if req.Status != "pending" {
		return fmt.Errorf("request cannot be rejected: current status is %s", req.Status)
	}

	// Обновление статуса на "rejected"
	if err := s.partnerRequestStorage.UpdateStatus(inn, "rejected"); err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

// Получение заявок по статусу
func (s *AdminManager) GetRequestsByStatus(status string) ([]*PartnerRequest, error) {
	return s.partnerRequestStorage.GetByStatus(status)
}

// Получение всех заявок
func (s *AdminManager) GetAllRequests() ([]*PartnerRequest, error) {
	return s.partnerRequestStorage.GetAll()
}

// Получение всех заявок в работе
func (s *AdminManager) GetPendingRequests() ([]*PartnerRequest, error) {
	return s.partnerRequestStorage.GetPending()
}

// Получение статуса заявки по ИНН
func (s *AdminManager) GetRequestStatus(inn string) (*PartnerRequest, error) {
	// Получение заявки по INN
	req, err := s.partnerRequestStorage.GetByINN(inn)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}
	if req == nil {
		return nil, fmt.Errorf("request with inn %s not found", inn)
	}

	return req, nil
}

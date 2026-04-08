package partners

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
type PartnersManager struct {
	userStorage           UserStorage
	partnerRequestStorage PartnerRequestStorage
	companyStorage        CompanyStorage
	emailSender           configPkg.EmailSender
	config                Config
}

// Создает новый экземпляр сервиса
func NewPartnersManager(
	userStorage UserStorage,
	partnerRequestStorage PartnerRequestStorage,
	companyStorage CompanyStorage,
	emailSender configPkg.EmailSender,
	config Config,
) *PartnersManager {
	return &PartnersManager{
		userStorage:           userStorage,
		partnerRequestStorage: partnerRequestStorage,
		companyStorage:        companyStorage,
		emailSender:           emailSender,
		config:                config,
	}
}

// Создание заявки партнера (доступно любому авторизованному пользователю)
func (s *PartnersManager) CreatePartnerRequest(userEmail string, req *PartnerRequestRequest) error {
	// Проверка на существование пользователя
	user, err := s.userStorage.GetByEmail(userEmail)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	isPartner, err := s.UserIsPartner(userEmail)
	if err != nil {
		return fmt.Errorf("failed to check if user is partner: %w", err)
	}

	// Проверка на наличие заявки с таким же ИНН
	existing, _ := s.companyStorage.GetByINN(req.INN)
	if existing != nil {
		return fmt.Errorf("company with this INN already exists")
	}

	existingRequest, err := s.partnerRequestStorage.GetByINN(req.INN)
	if err != nil {
		return fmt.Errorf("failed to check existing request: %w", err)
	}
	if existingRequest != nil {
		if existingRequest.Status != "rejected" {
			return fmt.Errorf("an active request with this INN already exists (status: %s)", existingRequest.Status)
		}
	}
	if isPartner.IsPartner {
		if existingRequest != nil {
			if existingRequest.Status != "rejected" {
				return fmt.Errorf("user is already a partner in a company or has an active request")
			}
		}
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

// Получение статуса заявки по ИНН
func (s *PartnersManager) GetRequestStatus(inn string) (*PartnerRequest, error) {
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

// Получение статуса заявки по email пользователя
func (s *PartnersManager) GetRequestStatusByEmail(email string) (*PartnerRequest, error) {
	req, err := s.partnerRequestStorage.GetByUserEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}
	if req == nil {
		return nil, fmt.Errorf("no request found for user: %s", email)
	}

	return req, nil
}

// Проверка что у пользователя есть организация
// Если у пользователя организация есть вернётся IsParnersUsers{IsPartner: true, Inn: "строка с инн"
// Если у пользователя организация нет вернётся IsParnersUsers{IsPartner: false, Inn: ничего
func (s *PartnersManager) UserIsPartner(email string) (IsPartnersUsers, error) {
	partUser, err := s.partnerRequestStorage.GetPartUserByEmail(email)
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

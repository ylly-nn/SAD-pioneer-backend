package branch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrDetailsNotFound = errors.New("details not found")
)

// BranchManager содержит бизнес-логику для управления связями филиалов и услуг.
type BranchManager struct {
	storage BranchStorage
}

// NewBranchManager создаёт новый BranchManager с заданным хранилищем.
func NewBranchManager(storage BranchStorage) *BranchManager {
	return &BranchManager{storage: storage}
}

// CreateBranchService создаёт новую запись в branch_services.
func (m *BranchManager) CreateBranchService(req CreateBranchServRequest) (*BranchServ, error) {
	if req.Branch == uuid.Nil {
		return nil, fmt.Errorf("branch is required")
	}
	if req.Service == uuid.Nil {
		return nil, fmt.Errorf("service is required")
	}
	if req.ServiceDetails == nil {
		return nil, fmt.Errorf("service_detalis is required")
	}

	// Строгая валидация: service_detalis должен быть объектом с ключами-строками и значениями-числами (минуты)
	var detailsMap map[string]int
	if err := json.Unmarshal(req.ServiceDetails, &detailsMap); err != nil {
		return nil, fmt.Errorf("service_detalis must be a JSON object with string keys and numeric values (minutes)")
	}

	// Проверка, что значения положительные
	for name, minutes := range detailsMap {
		if minutes <= 0 {
			return nil, fmt.Errorf("duration for '%s' must be positive (minutes)", name)
		}
	}

	if len(detailsMap) == 0 {
		return nil, fmt.Errorf("service_detalis cannot be empty")
	}

	bs := BranchServ{
		Branch:         req.Branch,
		Service:        req.Service,
		ServiceDetails: req.ServiceDetails,
		BusyTime:       nil,
	}
	return m.storage.CreateBranchServ(bs)
}

// CreateBranch создаёт новый филиал с проверкой бизнес-правил.
func (m *BranchManager) CreateBranch(req CreateBranchRequest) (*Branch, error) {
	// Проверка обязательных полей
	if req.City == "" {
		return nil, fmt.Errorf("city is required")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}
	if req.Inn == "" {
		return nil, fmt.Errorf("inn_company is required")
	}

	// ИНН должен быть ровно 12 цифр
	if len(req.Inn) != 12 {
		return nil, fmt.Errorf("inn_company must be exactly 12 characters")
	}
	for _, c := range req.Inn {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("inn_company must contain only digits")
		}
	}

	// Проверка времени
	if req.OpenTime.IsZero() {
		return nil, fmt.Errorf("open_time is required")
	}
	if req.CloseTime.IsZero() {
		return nil, fmt.Errorf("close_time is required")
	}
	if !req.OpenTime.Before(req.CloseTime) {
		return nil, fmt.Errorf("open_time must be before close_time")
	}

	branch := Branch{
		City:      req.City,
		Address:   req.Address,
		Inn:       req.Inn,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
	}

	created, err := m.storage.CreateBranch(branch)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}
	return created, nil
}

func (m *BranchManager) GetBranchByCityServ(city string, serviceID string) ([]*BrancByCityServ, error) {

	return m.storage.GetBranchByCityServ(city, serviceID)
}

func (m *BranchManager) GetServiceDetails(branchServID uuid.UUID) ([]*ServResponse, error) {
	details, price, err := m.storage.GetServiceDetails(branchServID)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[string]float32, len(price))
	for _, p := range price {
		priceMap[p.Detail] = p.Price
	}

	var result []*ServResponse
	for _, d := range details {
		// Ищем цену для текущей детали
		pr, ok := priceMap[d.Detail]
		if !ok {
			return nil, fmt.Errorf("price not found for detail: %s", d.Detail)
		}
		result = append(result, &ServResponse{
			Detail:   d.Detail,
			Duration: d.Duration,
			Price:    pr,
		})
	}

	if result == nil {
		return nil, ErrDetailsNotFound
	}

	return result, nil
}

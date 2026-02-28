package branch

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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
	return m.storage.Create(bs)
}

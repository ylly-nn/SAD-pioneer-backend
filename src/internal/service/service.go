package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	// ErrEmptyName возвращается, если название услуги пустое.
	ErrEmptyName = errors.New("service name cannot be empty")
)

// ServiceManager содержит бизнес-логику для работы с услугами.
type ServiceManager struct {
	storage ServiceStorage
}

// NewServiceManager создаёт новый экземпляр ServiceManager.
func NewServiceManager(storage ServiceStorage) *ServiceManager {
	return &ServiceManager{storage: storage}
}

// GetAllServices возвращает список всех услуг.
func (m *ServiceManager) GetAllServices() ([]Service, error) {
	services, err := m.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all services: %w", err)
	}
	return services, nil
}

// CreateService создаёт новую услугу.
func (m *ServiceManager) CreateService(name string) (*Service, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return m.storage.Create(name)
}

// DeleteService удаляет услугу по идентификатору.
// Если услуга не найдена, возвращает ошибку ErrServiceNotFound.
func (m *ServiceManager) DeleteService(id uuid.UUID) error {
	return m.storage.Delete(id)
}

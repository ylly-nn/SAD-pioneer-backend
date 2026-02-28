package client

import (
	"errors"
)

var (
	ErrEmptyEmail = errors.New("email cannot be empty")
)

// ClientManager содержит бизнес-логику для работы с владельцами тс (клиентами)
type ClientManager struct {
	storage ClientStorage
}

// NewClientManager создаёт новый экземпляр ClientManager
func NewClientManager(storage ClientStorage) *ClientManager {
	return &ClientManager{storage: storage}
}

// CreateClient создаёт нового клиента.
func (m *ClientManager) CreateClient(email string) (*Client, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}
	return m.storage.Create(email)
}

// UpdateCity обновляет город клиента.
// Если email пуст, возвращает ErrEmptyEmail.
// Ошибки storage пробрасываются наружу.
func (m *ClientManager) UpdateCity(email string, city string) error {
	if email == "" {
		return ErrEmptyEmail
	}
	return m.storage.UpdateCity(email, city)
}

// GetCityByEmail возвращает город клиента.
// Если email пуст, возвращает ErrEmptyEmail.
// Если клиент не найден, возвращает ErrClientNotFound.
func (m *ClientManager) GetCityByEmail(email string) (string, error) {
	if email == "" {
		return "", ErrEmptyEmail
	}
	return m.storage.GetCityByEmail(email)
}

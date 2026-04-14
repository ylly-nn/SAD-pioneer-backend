package client

import (
	"errors"
	"regexp"
	"strings"

	"src/internal/city"
)

var (
	ErrEmptyEmail  = errors.New("email cannot be empty")
	ErrEmptyCity   = errors.New("city cannot be empty")
	ErrInvalidCity = errors.New("city is not in the list of Russian cities")
)

var hyphenSpaces = regexp.MustCompile(`\s*-\s*`)

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
func (m *ClientManager) UpdateCity(email string, cityName string) error {
	if email == "" {
		return ErrEmptyEmail
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

	return m.storage.UpdateCity(email, canonicalCity)
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

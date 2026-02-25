package service

import (
	"github.com/google/uuid"
)

// Service представляет данные общей услуги
// Соответствует таблице services в базе данных
type Service struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// CreateServiceRequest описывает тело запроса для создания новой услуги.
// Используется в эндпоинте POST /services.
type CreateServiceRequest struct {
	Name string `json:"name"`
}

// ServiceResponse представляет данные услуги, возвращаемые клиенту.
// Используется в ответах на GET /services
type ServiceResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

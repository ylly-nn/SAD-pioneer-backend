package client

import "github.com/google/uuid"

// Client представляет данные клиента из таблицы ts_users.
// Используется для внутренней передачи данных.
type Client struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	City  *string   `json:"city,omitempty"`
}

// ClientResponse содержит данные для ответа на GET /client.
type ClientResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	City  *string   `json:"city,omitempty"`
}

// CreateClientRequest содержит данные для создания клиента через POST /client.
type CreateClientRequest struct {
	Email string `json:"email"`
}

package client

import "github.com/google/uuid"

// Соответствует таблице ts_users
// Для внутренней передачи данных
type Client struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	City  *string   `json:"city,omitempty"`
}

// Структура для get /client
type ClientResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	City  *string   `json:"city,omitempty"`
}

// Структура для POST /client
type CreateClientRequest struct {
	Email string `json:"email"`
}

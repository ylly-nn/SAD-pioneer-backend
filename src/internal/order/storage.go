package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"src/internal/db"
)

// OrderStorage определяет методы для работы с услугами в базе данных.
type OrderStorage interface {
	// Create(users uuid.UUID, ServiceByBranch uuid.UUID, Date time.Time,
	// 	StartTime time.Time, OrderDetails json.RawMessage) (*Order, error)

	GetAll() ([]*Order, error)

	//GetByClient(email string) ([]*Order, error)

	//GetByBranch(branchID uuid.UUID) ([]*Order, error)

	//GetByCompany

	//Delete?

	//Update?
}

// реализует OrderStorage для PostgreSQL.
type PostgresOrderStorage struct {
	*db.Storage
}

func NewPostrgesOrderStorage(sqlDB *sql.DB) *PostgresOrderStorage {
	return &PostgresOrderStorage{Storage: db.NewStorage(sqlDB)}
}

// GetAll возвращает список всех заказов из таблицы orders.
func (s *PostgresOrderStorage) GetAll() ([]*Order, error) {
	rows, err := s.DB.Query(`
        SELECT id, users, service_by_branch, date, start_time, order_details
        FROM orders
        ORDER BY date, start_time
    `)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var ord Order
		var date sql.NullTime
		var startTimeStr sql.NullString
		var orderDetails []byte

		err := rows.Scan(
			&ord.ID,
			&ord.Users,
			&ord.ServiceByBranch,
			&date,
			&startTimeStr,
			&orderDetails,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if date.Valid {
			ord.Date = &date.Time
		}

		// Преобразуем start_time из строки в time.Time
		if startTimeStr.Valid {
			t, err := time.Parse("15:04:05-07", startTimeStr.String)
			if err != nil {
				t, err = time.Parse("15:04:05Z07:00", startTimeStr.String)
				if err != nil {
					return nil, fmt.Errorf("failed to parse start_time %q: %w", startTimeStr.String, err)
				}
			}
			ord.StartTime = &t
		}

		if len(orderDetails) > 0 {
			ord.OrderDetails = json.RawMessage(orderDetails)
		}

		orders = append(orders, &ord)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return orders, nil
}

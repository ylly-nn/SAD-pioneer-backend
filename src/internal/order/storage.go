package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"src/internal/db"

	"github.com/google/uuid"
)

// OrderStorage определяет методы для работы с услугами в базе данных.
type OrderStorage interface {
	GetFullAllOrders() ([]*FullOrder, error)

	GetByCLient(email string) ([]*FullOrder, error)

	GetByCompany(inn string) ([]*FullOrder, error)

	Create(order Order) (*Order, error)
}

// реализует OrderStorage для PostgreSQL.
type PostgresOrderStorage struct {
	*db.Storage
}

func NewPostrgesOrderStorage(sqlDB *sql.DB) *PostgresOrderStorage {
	return &PostgresOrderStorage{Storage: db.NewStorage(sqlDB)}
}

// Create добавляет новый заказ в таблицу orders.
// Все поля (users, service_by_branch, date, start_time, order_details) обязательны.
// Возвращает созданный заказ с заполненным ID.
func (s *PostgresOrderStorage) Create(order Order) (*Order, error) {
	if order.Users == uuid.Nil {
		return nil, fmt.Errorf("users is required")
	}
	if order.ServiceByBranch == uuid.Nil {
		return nil, fmt.Errorf("service_by_branch is required")
	}
	if order.Date == nil || order.Date.IsZero() {
		return nil, fmt.Errorf("date is required")
	}
	if order.StartTime == nil || order.StartTime.IsZero() {
		return nil, fmt.Errorf("start_time is required")
	}
	if len(order.OrderDetails) == 0 {
		return nil, fmt.Errorf("order_details is required")
	}

	var id uuid.UUID
	err := s.DB.QueryRow(`
		INSERT INTO orders (users, service_by_branch, date, start_time, order_details)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, order.Users, order.ServiceByBranch, *order.Date, *order.StartTime, order.OrderDetails).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	createdOrder := &Order{
		ID:              id,
		Users:           order.Users,
		ServiceByBranch: order.ServiceByBranch,
		Date:            order.Date,
		StartTime:       order.StartTime,
		OrderDetails:    order.OrderDetails,
	}
	return createdOrder, nil
}

// Выводит полную информацию о заказе для определённой оргвнизации
func (s *PostgresOrderStorage) GetByCompany(inn string) ([]*FullOrder, error) {
	rows, err := s.DB.Query(`
        SELECT o.id, o.users, ts.email, o.service_by_branch, b.inn_company, c.org_short_name, b.city, b.address, s.name, o.date, o.start_time, o.order_details 
        FROM orders o
        JOIN ts_users ts ON o.users = ts.id
        JOIN branch_services bs ON o.service_by_branch = bs.id
        JOIN services s on bs.service = s.id
        JOIN branches b on bs.branch=b.id
        JOIN companies c on b.inn_company = c.inn
        WHERE c.inn = $1;

    `, inn)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var orders []*FullOrder
	for rows.Next() {
		var ord FullOrder
		var date sql.NullTime
		var startTimeStr sql.NullString
		var orderDetails []byte

		err := rows.Scan(
			&ord.ID,
			&ord.Users,
			&ord.Email,
			&ord.ServiceByBranch,
			&ord.InnCompany,
			&ord.NameCompany,
			&ord.City,
			&ord.Address,
			&ord.Service,
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

// Выводит полную информацию о заказе для определённого клиента
func (s *PostgresOrderStorage) GetByCLient(email string) ([]*FullOrder, error) {
	rows, err := s.DB.Query(`
        SELECT o.id, o.users, ts.email, o.service_by_branch, b.inn_company, c.org_short_name, b.city, b.address, s.name, o.date, o.start_time, o.order_details 
        FROM orders o
        JOIN ts_users ts ON o.users = ts.id
        JOIN branch_services bs ON o.service_by_branch = bs.id
        JOIN services s ON bs.service = s.id
        JOIN branches b ON bs.branch = b.id
        JOIN companies c ON b.inn_company = c.inn
        WHERE ts.email = $1;
    `, email)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var orders []*FullOrder
	for rows.Next() {
		var ord FullOrder
		var date sql.NullTime
		var startTimeStr sql.NullString
		var orderDetails []byte

		err := rows.Scan(
			&ord.ID,
			&ord.Users,
			&ord.Email,
			&ord.ServiceByBranch,
			&ord.InnCompany,
			&ord.NameCompany,
			&ord.City,
			&ord.Address,
			&ord.Service,
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

// выводит инфлмации о заказе с клентом, сервисом, филиалом и команией
func (s *PostgresOrderStorage) GetFullAllOrders() ([]*FullOrder, error) {
	rows, err := s.DB.Query(`
	SELECT o.id, o.users, ts.email, o.service_by_branch, b.inn_company, c.org_short_name, 
	b.city, b.address, s.name, o.date, o.start_time, o.order_details 
        FROM orders o
        JOIN ts_users ts ON o.users = ts.id
        JOIN branch_services bs ON o.service_by_branch = bs.id
        JOIN services s on bs.service = s.id
        JOIN branches b on bs.branch=b.id
        JOIN companies c on b.inn_company = c.inn;
	`)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var orders []*FullOrder
	for rows.Next() {
		var ord FullOrder
		var date sql.NullTime
		var startTimeStr sql.NullString
		var orderDetails []byte

		err := rows.Scan(
			&ord.ID,
			&ord.Users,
			&ord.Email,
			&ord.ServiceByBranch,
			&ord.InnCompany,
			&ord.NameCompany,
			&ord.City,
			&ord.Address,
			&ord.Service,
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

		// Обработка времени начала (start_time) – парсинг из строки с учётом временной зоны
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

		// Обработка JSON-поля order_details
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

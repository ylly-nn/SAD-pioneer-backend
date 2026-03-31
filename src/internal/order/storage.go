package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"src/internal/db"

	"github.com/google/uuid"
)

// OrderStorage определяет методы для работы с услугами в базе данных.
type OrderStorage interface {
	GetByClient(email string) ([]*FullOrder, error)

	Create(order Order) (*Order, error)

	GetOpenCloseTime(branch_id uuid.UUID) (*OpenCloseBranch, error)

	GetBisyTimeByDate(branch_id uuid.UUID, date time.Time) ([]*BusyTime, error)

	GetBranchIDByBranchServ(branchServID uuid.UUID) (uuid.UUID, error)

	GetDetailsByBranchServ(branchServID uuid.UUID) ([]*ServiceDuration, []*ServPrice, error)
	//GetFullAllOrders() ([]*FullOrder, error)
	//GetByCompany(inn string) ([]*FullOrder, error)
}

var (
	ErrBranchServiceNotFound = errors.New("branch service not found")
	ErrOrdersNotFound        = errors.New("orders not found")
)

// реализует OrderStorage для PostgreSQL.
type PostgresOrderStorage struct {
	*db.Storage
}

func NewPostrgesOrderStorage(sqlDB *sql.DB) *PostgresOrderStorage {
	return &PostgresOrderStorage{Storage: db.NewStorage(sqlDB)}
}

// полуение деталей услуги и их стоимости в виде  []*ServiceDuration и  []*ServPrice
func (s *PostgresOrderStorage) GetDetailsByBranchServ(branchServID uuid.UUID) ([]*ServiceDuration, []*ServPrice, error) {
	var detailsJSON json.RawMessage
	var priceJSON json.RawMessage
	err := s.DB.QueryRow(`SELECT service_detalis, price FROM branch_services WHERE id = $1`, branchServID).Scan(&detailsJSON, &priceJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("branch service not found")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrBranchServiceNotFound
		}
		return nil, nil, fmt.Errorf("query failed: %w", err)
	}

	// Если JSON пустой возвращаем пустой срез
	if len(detailsJSON) == 0 || string(detailsJSON) == "null" {
		return []*ServiceDuration{}, []*ServPrice{}, nil
	}
	var detailsMap map[string]int
	if len(detailsJSON) == 0 || string(detailsJSON) == "null" {
		detailsMap = make(map[string]int)
	} else {
		if err := json.Unmarshal(detailsJSON, &detailsMap); err != nil {
			return nil, nil, fmt.Errorf("failed to parse service details: %w", err)
		}
	}

	var priceMap map[string]float32
	if len(priceJSON) == 0 || string(priceJSON) == "null" {
		priceMap = make(map[string]float32)
	} else {
		if err := json.Unmarshal(priceJSON, &priceMap); err != nil {
			return nil, nil, fmt.Errorf("failed to parse service price: %w", err)
		}
	}

	detailsStruct := make([]*ServiceDuration, 0, len(detailsMap))
	for detail, duration := range detailsMap {
		detailsStruct = append(detailsStruct, &ServiceDuration{
			Detail:   detail,
			Duration: duration,
		})
	}

	priceStruct := make([]*ServPrice, 0, len(priceMap))
	for detail, price := range priceMap {
		priceStruct = append(priceStruct, &ServPrice{
			Detail: detail,
			Price:  price,
		})
	}

	return detailsStruct, priceStruct, nil
}

// Create добавляет новый заказ в таблицу orders
// Возвращает созданный заказ с заполненным ID.
func (s *PostgresOrderStorage) Create(order Order) (*Order, error) {
	order.Status = OrderStatusCreate
	var id uuid.UUID
	err := s.DB.QueryRow(`
		INSERT INTO orders (users, service_by_branch, start_moment, end_moment, order_details, status, price, sum)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, order.Users, order.ServiceByBranch, order.StartMoment, order.EndMoment, order.OrderDetails, order.Status, order.Price, order.Sum).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	createdOrder := &Order{
		ID:              id,
		Users:           order.Users,
		ServiceByBranch: order.ServiceByBranch,
		StartMoment:     order.StartMoment,
		EndMoment:       order.EndMoment,
		Status:          order.Status,
		OrderDetails:    order.OrderDetails,
		Price:           order.Price,
		Sum:             order.Sum,
	}
	return createdOrder, nil
}

// GetOpenCloseTime возвращает время открытия и закрытия филиала по его ID
// Возвращает структуру OpenCloseBranch или ошибку, если филиал не найден
func (s *PostgresOrderStorage) GetOpenCloseTime(branchID uuid.UUID) (*OpenCloseBranch, error) {
	var openTimeStr, closeTimeStr sql.NullString
	err := s.DB.QueryRow(`
        SELECT open_time, close_time
        FROM branches
        WHERE id = $1
    `, branchID).Scan(&openTimeStr, &closeTimeStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("branch with id %s not found", branchID)
		}
		return nil, fmt.Errorf("failed to query open/close time: %w", err)
	}
	if !openTimeStr.Valid {
		return nil, fmt.Errorf("open_time is null for branch %s", branchID)
	}
	if !closeTimeStr.Valid {
		return nil, fmt.Errorf("close_time is null for branch %s", branchID)
	}

	openTime, err := time.Parse("15:04:05-07:00", openTimeStr.String)
	if err != nil {
		openTime, err = time.Parse("15:04:05-07", openTimeStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse open_time %q: %w", openTimeStr.String, err)
		}
	}
	closeTime, err := time.Parse("15:04:05-07:00", closeTimeStr.String)
	if err != nil {
		closeTime, err = time.Parse("15:04:05-07", closeTimeStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close_time %q: %w", closeTimeStr.String, err)
		}
	}

	return &OpenCloseBranch{
		OpenTimeBranch:  openTime,
		CloseTimeBranch: closeTime,
	}, nil
}

// GetBisyTimeByDate возвращает список занятых интервалов для филиала на указанную дату
// Возвращает срез BusyTime или ошибку при выполнении запроса
func (s *PostgresOrderStorage) GetBisyTimeByDate(branchID uuid.UUID, date time.Time) ([]*BusyTime, error) {
	rows, err := s.DB.Query(`
        SELECT o.start_moment, o.end_moment
        FROM orders o
        JOIN branch_services bs ON o.service_by_branch = bs.id
        WHERE bs.branch = $1 AND o.start_moment::date = $2::date AND status != 'reject'
    `, branchID, date)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var intervals []*BusyTime
	for rows.Next() {
		var start, end sql.NullTime
		if err := rows.Scan(&start, &end); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if !start.Valid {
			return nil, fmt.Errorf("unexpected null start_moment")
		}

		ot := &BusyTime{
			StartMoment: start.Time,
		}
		if end.Valid {
			ot.EndMoment = end.Time
		}

		intervals = append(intervals, ot)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return intervals, nil
}

// GetBranchIDByBranchServ возвращает UUID филиала, связанного с указанной услугой филиала.
func (s *PostgresOrderStorage) GetBranchIDByBranchServ(branchServID uuid.UUID) (uuid.UUID, error) {
	var branchID uuid.UUID
	err := s.DB.QueryRow(`
        SELECT branch FROM branch_services WHERE id = $1
    `, branchServID).Scan(&branchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("branch_service with id %s not found", branchServID)
		}
		return uuid.Nil, fmt.Errorf("failed to query branch id: %w", err)
	}
	return branchID, nil
}

// Выводит полную информацию о заказе для определённого клиента
func (s *PostgresOrderStorage) GetByClient(email string) ([]*FullOrder, error) {
	rows, err := s.DB.Query(`
        SELECT o.id, ts.email, o.service_by_branch, b.inn_company, c.org_short_name, b.city, b.address, s.name, o.start_moment, o.end_moment, o.order_details, o.status, o.price, o.sum
        FROM orders o
        JOIN ts_users ts ON o.users = ts.email
        JOIN branch_services bs ON o.service_by_branch = bs.id
        JOIN services s ON bs.service = s.id
        JOIN branches b ON bs.branch = b.id
        JOIN companies c ON b.inn_company = c.inn
        WHERE ts.email = $1
    `, email)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var orders []*FullOrder
	for rows.Next() {
		var ord FullOrder
		var endMoment sql.NullTime
		var orderDetailsRaw []byte
		var priceRaw []byte
		var sum float32

		err := rows.Scan(
			&ord.ID,
			&ord.Email,
			&ord.ServiceByBranch,
			&ord.InnCompany,
			&ord.NameCompany,
			&ord.City,
			&ord.Address,
			&ord.Service,
			&ord.StartMoment,
			&endMoment,
			&orderDetailsRaw,
			&ord.Status,
			&priceRaw,
			&sum,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if endMoment.Valid {
			ord.EndMoment = &endMoment.Time
		}

		// Десериализация order_details (хранится как map[string]int)
		if len(orderDetailsRaw) > 0 {
			var detailsMap map[string]int
			if err := json.Unmarshal(orderDetailsRaw, &detailsMap); err != nil {
				return nil, fmt.Errorf("failed to unmarshal order_details: %w", err)
			}
			detailsSlice := make([]ServiceDuration, 0, len(detailsMap))
			for detail, duration := range detailsMap {
				detailsSlice = append(detailsSlice, ServiceDuration{
					Detail:   detail,
					Duration: duration,
				})
			}
			ord.OrderDetails = detailsSlice
		} else {
			ord.OrderDetails = []ServiceDuration{}
		}

		// Десериализация price (хранится как map[string]float32)
		if len(priceRaw) > 0 {
			var priceMap map[string]float32
			if err := json.Unmarshal(priceRaw, &priceMap); err != nil {
				return nil, fmt.Errorf("failed to unmarshal price: %w", err)
			}
			priceSlice := make([]ServPrice, 0, len(priceMap))
			for detail, priceVal := range priceMap {
				priceSlice = append(priceSlice, ServPrice{
					Detail: detail,
					Price:  priceVal,
				})
			}
			ord.Price = priceSlice
		} else {
			ord.Price = []ServPrice{}
		}

		ord.Sum = sum
		orders = append(orders, &ord)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}
	return orders, nil
}

// // выводит инфлмации о заказе с клентом, сервисом, филиалом и команией
// func (s *PostgresOrderStorage) GetFullAllOrders() ([]*FullOrder, error) {
// 	rows, err := s.DB.Query(`
//         SELECT o.id, ts.email, o.service_by_branch, b.inn_company, c.org_short_name,
//                b.city, b.address, s.name, o.start_moment, o.end_moment, o.order_details
//         FROM orders o
//         JOIN ts_users ts ON o.users = ts.email
//         JOIN branch_services bs ON o.service_by_branch = bs.id
//         JOIN services s ON bs.service = s.id
//         JOIN branches b ON bs.branch = b.id
//         JOIN companies c ON b.inn_company = c.inn
//     `)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var orders []*FullOrder
// 	for rows.Next() {
// 		var ord FullOrder
// 		var endMoment sql.NullTime
// 		var orderDetails []byte

// 		err := rows.Scan(
// 			&ord.ID,
// 			&ord.Email,
// 			&ord.ServiceByBranch,
// 			&ord.InnCompany,
// 			&ord.NameCompany,
// 			&ord.City,
// 			&ord.Address,
// 			&ord.Service,
// 			&ord.StartMoment,
// 			&endMoment,
// 			&orderDetails,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}

// 		if endMoment.Valid {
// 			ord.EndMoment = &endMoment.Time
// 		}

// 		if len(orderDetails) > 0 {
// 			ord.OrderDetails = json.RawMessage(orderDetails)
// 		}

// 		orders = append(orders, &ord)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows iteration error: %w", err)
// 	}
// 	return orders, nil
// }

// // Выводит полную информацию о заказе для определённой оргвнизации
// func (s *PostgresOrderStorage) GetByCompany(inn string) ([]*FullOrder, error) {
// 	rows, err := s.DB.Query(`
//         SELECT o.id, ts.email, o.service_by_branch, b.inn_company, c.org_short_name, b.city, b.address, s.name, o.start_moment, o.end_moment, o.order_details
//         FROM orders o
//         JOIN ts_users ts ON o.users = ts.email
//         JOIN branch_services bs ON o.service_by_branch = bs.id
//         JOIN services s ON bs.service = s.id
//         JOIN branches b ON bs.branch = b.id
//         JOIN companies c ON b.inn_company = c.inn
//         WHERE c.inn = $1;
//     `, inn)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var orders []*FullOrder
// 	for rows.Next() {
// 		var ord FullOrder
// 		var endMoment sql.NullTime
// 		var orderDetails []byte

// 		err := rows.Scan(
// 			&ord.ID,
// 			&ord.Email,
// 			&ord.ServiceByBranch,
// 			&ord.InnCompany,
// 			&ord.NameCompany,
// 			&ord.City,
// 			&ord.Address,
// 			&ord.Service,
// 			&ord.StartMoment,
// 			&endMoment,
// 			&orderDetails,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}

// 		if endMoment.Valid {
// 			ord.EndMoment = &endMoment.Time
// 		}

// 		if len(orderDetails) > 0 {
// 			ord.OrderDetails = json.RawMessage(orderDetails)
// 		}

// 		orders = append(orders, &ord)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows iteration error: %w", err)
// 	}
// 	return orders, nil
// }

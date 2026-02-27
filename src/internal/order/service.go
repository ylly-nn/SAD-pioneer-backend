package order

// содержит бизнес-логику для работы с услугами.
type OrderManager struct {
	storage OrderStorage
}

// создаёт новый экземпляр OrderManager.
func NewOrderManager(storage OrderStorage) *OrderManager {
	return &OrderManager{storage: storage}
}

// GetAllOrders возвращает список всех заказов.
func (m *OrderManager) GetAllOrders() ([]*Order, error) {
	return m.storage.GetAll()
}

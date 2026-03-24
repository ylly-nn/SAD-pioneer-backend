package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyEmail = errors.New("email cannot be empty")
	ErrInnLen     = errors.New("inn length must be 12 characters")
)

// содержит бизнес-логику для работы с услугами.
type OrderManager struct {
	storage OrderStorage
}

// создаёт новый экземпляр OrderManager.
func NewOrderManager(storage OrderStorage) *OrderManager {
	return &OrderManager{storage: storage}
}

// Create создаёт новый заказ после проверки доступности выбранного времени.
func (m *OrderManager) Create(email string, req CreateOrderRequest) (*Order, error) { // время окончания
	// Приводим время к UTC для согласованности с БД
	req.StartMoment = req.StartMoment.UTC()

	// Проверка обязательных полей
	if email == "" {
		return nil, errors.New("users is required")
	}
	if req.ServiceByBranch == uuid.Nil {
		return nil, errors.New("service_by_branch is required")
	}
	if req.StartMoment.IsZero() {
		return nil, errors.New("start moment is required")
	}
	if len(req.OrderDetails) == 0 {
		return nil, errors.New("order_details is required")
	}

	// Валидация order_details: должен быть объектом { "услуга": минуты }
	var details map[string]int
	if err := json.Unmarshal(req.OrderDetails, &details); err != nil {
		return nil, errors.New("order_details must be a JSON object with string keys and numeric values (minutes)")
	}
	if len(details) == 0 {
		return nil, errors.New("order_details cannot be empty")
	}

	// Подсчёт общей длительности
	totalMinutes := 0
	for name, minutes := range details {
		if minutes <= 0 {
			return nil, fmt.Errorf("duration for '%s' must be positive (minutes)", name)
		}
		totalMinutes += minutes
	}

	// Получение ID филиала по услуге филиала
	branchID, err := m.storage.GetBranchIDByBranchServ(req.ServiceByBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch for service: %w", err)
	}

	// Проверка доступен ли слот запрошенной длительности
	slots, err := m.GetFreeTimeForDay(branchID, req.StartMoment, totalMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to check free time: %w", err)
	}

	startUTC := req.StartMoment.UTC()
	valid := false
	for _, slot := range slots {
		if slot.UTC().Equal(startUTC) {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("start moment is not available for the requested duration")
	}

	endTime := req.StartMoment.Add(time.Duration(totalMinutes) * time.Minute)

	order := Order{
		Users:           email,
		ServiceByBranch: req.ServiceByBranch,
		StartMoment:     req.StartMoment,
		EndMoment:       &endTime,
		OrderDetails:    req.OrderDetails,
	}
	return m.storage.Create(order)
}

// GetFreeTimeForWeek возвращает свободные слоты с шагом 15 минут
// для указанного филиала на день
func (m *OrderManager) GetFreeTimeForDay(branchID uuid.UUID, day time.Time, duration int) ([]time.Time, error) {
	openClose, err := m.storage.GetOpenCloseTime(branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get open/close time: %w", err)
	}

	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

	busy, err := m.storage.GetBisyTimeByDate(branchID, dayStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get busy time for %s: %w", dayStart.Format("2006-01-02"), err)
	}
	if busy == nil {
		busy = []*BusyTime{}
	}

	sort.Slice(busy, func(i, j int) bool {
		return busy[i].StartMoment.Before(busy[j].StartMoment)
	})

	openTime := time.Date(dayStart.Year(), dayStart.Month(), dayStart.Day(),
		openClose.OpenTimeBranch.Hour(), openClose.OpenTimeBranch.Minute(),
		openClose.OpenTimeBranch.Second(), 0, openClose.OpenTimeBranch.Location())
	closeTime := time.Date(dayStart.Year(), dayStart.Month(), dayStart.Day(),
		openClose.CloseTimeBranch.Hour(), openClose.CloseTimeBranch.Minute(),
		openClose.CloseTimeBranch.Second(), 0, openClose.CloseTimeBranch.Location())

	slots := computeFreeSlots(openTime, closeTime, busy, duration)
	return slots, nil
}

// GetFreeTimeForWeek возвращает свободные слоты с шагом 15 минут
// для указанного филиала на неделю, начиная с startDate.
func (m *OrderManager) GetFreeTimeForWeek(branchID uuid.UUID, startDate time.Time, duration int) ([]DailySlots, error) {
	openClose, err := m.storage.GetOpenCloseTime(branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get open/close time: %w", err)
	}

	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	var weekFree []DailySlots

	for i := range 7 {
		day := start.AddDate(0, 0, i)

		busy, err := m.storage.GetBisyTimeByDate(branchID, day)
		if err != nil {
			return nil, fmt.Errorf("failed to get busy time for %s: %w", day.Format("2006-01-02"), err)
		}
		if busy == nil {
			busy = []*BusyTime{}
		}

		sort.Slice(busy, func(i, j int) bool {
			return busy[i].StartMoment.Before(busy[j].StartMoment)
		})

		// Формируем время открытия и закрытия для текущего дня
		openTime := time.Date(day.Year(), day.Month(), day.Day(),
			openClose.OpenTimeBranch.Hour(), openClose.OpenTimeBranch.Minute(),
			openClose.OpenTimeBranch.Second(), 0, openClose.OpenTimeBranch.Location())
		closeTime := time.Date(day.Year(), day.Month(), day.Day(),
			openClose.CloseTimeBranch.Hour(), openClose.CloseTimeBranch.Minute(),
			openClose.CloseTimeBranch.Second(), 0, openClose.CloseTimeBranch.Location())

		slots := computeFreeSlots(openTime, closeTime, busy, duration)

		weekFree = append(weekFree, DailySlots{
			Date:      day,
			Intervals: slots,
		})
	}
	return weekFree, nil
}

// computeFreeIntervals вычисляет свободные промежутки между open и close
// с учётом занятых интервалов (busy). Предполагается, что busy отсортированы.
func computeFreeIntervals(open, close time.Time, busy []*BusyTime) []*BusyTime {
	var free []*BusyTime
	current := open

	for _, b := range busy {
		// Если начало занятого интервала позже текущего свободного времени
		if b.StartMoment.After(current) {
			free = append(free, &BusyTime{
				StartMoment: current,
				EndMoment:   b.StartMoment,
			})
		}
		// Передвигаем текущее время на конец занятого интервала (если он позже)
		if b.EndMoment.After(current) {
			current = b.EndMoment
		}
		// Если занятый интервал полностью внутри предыдущего свободного, current не меняется
	}

	// Добавляем остаток до закрытия
	if current.Before(close) {
		free = append(free, &BusyTime{
			StartMoment: current,
			EndMoment:   close,
		})
	}
	return free
}

// computeFreeSlots вычисляет все времена начала слотов длительностью duration минут,
// которые полностью помещаются в свободных промежутках между open и close с учётом busy.
func computeFreeSlots(open, close time.Time, busy []*BusyTime, duration int) []time.Time {
	var slots []time.Time
	current := open

	for _, b := range busy {
		if b.StartMoment.After(current) {
			// свободный промежуток [current, b.StartMoment)
			slots = append(slots, generateSlots(current, b.StartMoment, duration)...)
		}
		if b.EndMoment.After(current) {
			current = b.EndMoment
		}
	}
	if current.Before(close) {
		slots = append(slots, generateSlots(current, close, duration)...)
	}
	return slots
}

// generateSlots генерирует все времена начала слотов длительностью duration минут,
// укладывающихся в интервал [start, end] с шагом 15 минут.
func generateSlots(start, end time.Time, duration int) []time.Time {
	var slots []time.Time
	step := 15 * time.Minute
	dur := time.Duration(duration) * time.Minute

	for t := start; !t.Add(dur).After(end); t = t.Add(step) {
		slots = append(slots, t)
	}
	return slots
}

// возвращает список всех заказов
func (m *OrderManager) GetFullAllOrders() ([]*FullOrder, error) {
	return m.storage.GetFullAllOrders()
}

// возвращает список заказов опредеоённого киента
func (m *OrderManager) GetByClient(email string) ([]*ClientOrderResponse, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}

	fullOrders, err := m.storage.GetByClient(email)
	if err != nil {
		return nil, err
	}
	var responses []*ClientOrderResponse
	for _, fo := range fullOrders {
		responses = append(responses, &ClientOrderResponse{
			ID:           fo.ID,
			NameCompany:  fo.NameCompany,
			City:         fo.City,
			Address:      fo.Address,
			Service:      fo.Service,
			StartMoment:  fo.StartMoment,
			EndMoment:    fo.EndMoment,
			Status:       fo.Status,
			OrderDetails: fo.OrderDetails,
		})
	}
	return responses, nil
}

func (m *OrderManager) GetByCompany(inn string) ([]*FullOrder, error) {
	if len(inn) != 12 {
		return nil, ErrInnLen
	}
	return m.storage.GetByCompany(inn)
}

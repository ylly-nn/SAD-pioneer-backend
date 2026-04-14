package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyEmail              = errors.New("email cannot be empty")
	ErrInnLen                  = errors.New("inn length must be 12 characters")
	ErrDetailNameIsEpmpty      = errors.New("detail name cannot be empty")
	ErrDetailNotFound          = errors.New("detail with this name not found")
	ErrDetailNotAvailable      = errors.New("detail  is not available for this service")
	ErrDateInPast              = errors.New("date cannot be in the past")
	ErrDateInFuture            = errors.New("date cannot be more than one year in the future")
	ErrBranchServIsEmpty       = errors.New("service_by_branch cannot be empty")
	ErrStartMomentIsEmpty      = errors.New("start moment cannot be empty")
	ErrOrderDetailsIsEmpty     = errors.New("order_details cannot be empty")
	ErrTimeInPast              = errors.New("time cannot be in the past")
	ErrTimeInFuture            = errors.New("time cannot be more than one year in the future")
	ErrStartMomemtNotAvailable = errors.New("start moment is not available for the requested details")
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

	if req.ServiceByBranch == uuid.Nil {
		return nil, ErrBranchServIsEmpty
	}
	if req.StartMoment.IsZero() {
		return nil, ErrStartMomentIsEmpty
	}
	if len(req.OrderDetails) == 0 {
		return nil, ErrOrderDetailsIsEmpty
	}

	for _, d := range req.OrderDetails {
		if d.Detail == "" {
			return nil, ErrOrderDetailsIsEmpty
		}
	}

	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
	if req.StartMoment.Before(nowUTC) {
		return nil, ErrTimeInPast
	}

	maxDateUTC := todayUTC.AddDate(1, 0, 0)
	if req.StartMoment.After(maxDateUTC) {
		return nil, ErrTimeInFuture
	}

	detailsDB, priceDB, err := m.storage.GetDetailsByBranchServ(req.ServiceByBranch)
	if err != nil {
		return nil, err
	}
	if len(detailsDB) == 0 || len(priceDB) == 0 {
		return nil, ErrDetailNotFound
	}

	priceMap := make(map[string]float32)
	for _, p := range priceDB {
		priceMap[p.Detail] = p.Price
	}
	var servDetails []ServDetails

	for _, d := range detailsDB {
		price, exists := priceMap[d.Detail]
		if !exists {
			// Если цены нет, деталь не включается в заказ
			continue
		}
		servDetails = append(servDetails, ServDetails{
			Detail:   d.Detail,
			Duration: d.Duration,
			Price:    price,
		})

	}

	if len(servDetails) == 0 {
		return nil, ErrDetailNotFound
	}

	servDetailsJSON, _ := json.Marshal(servDetails)
	log.Printf("ServDetails: %s", servDetailsJSON)

	detailMap := make(map[string]ServDetails)
	for _, sd := range servDetails {
		detailMap[sd.Detail] = sd
	}

	// Формируем map деталей запроса с длительностями из БД
	detailsReq := make(map[string]int)
	priceReq := make(map[string]float32)
	var totalPrice float32 = 0

	for _, d := range req.OrderDetails {
		if d.Detail == "" {
			return nil, ErrDetailNameIsEpmpty
		}
		sd, ok := detailMap[d.Detail]
		if !ok {
			return nil, ErrDetailNotAvailable
		}
		// Берем длительность из базы
		detailsReq[d.Detail] = sd.Duration

		// Берем цену из ранее созданного priceMap
		price, exists := priceMap[d.Detail]
		if !exists {
			return nil, fmt.Errorf("price for detail '%s' not found", d.Detail)
		}
		priceReq[d.Detail] = price
		totalPrice += price
	}

	// Подсчёт общей длительности
	totalMinutes := 0
	for name, minutes := range detailsReq {
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
		return nil, ErrStartMomemtNotAvailable
	}

	endTime := req.StartMoment.Add(time.Duration(totalMinutes) * time.Minute)

	orderDetailsJSON, err := json.Marshal(detailsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order details: %w", err)
	}

	orderPriceJSON, err := json.Marshal(priceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order prices: %w", err)
	}

	order := Order{
		Users:           email,
		ServiceByBranch: req.ServiceByBranch,
		StartMoment:     req.StartMoment,
		EndMoment:       &endTime,
		OrderDetails:    orderDetailsJSON,
		Price:           orderPriceJSON,
		Sum:             totalPrice,
	}
	orderRes, err := m.storage.Create(order)
	if err != nil {
		return nil, err
	}

	return orderRes, err
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

	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
	if startDate.Before(todayUTC) {
		return []DailySlots{}, ErrDateInPast
	}

	maxDateUTC := todayUTC.AddDate(1, 0, 0)
	if startDate.After(maxDateUTC) {
		return []DailySlots{}, ErrDateInFuture
	}

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

		slotsTime := computeFreeSlots(openTime, closeTime, busy, duration)
		slots := make([]UTCTime, len(slotsTime))
		for i, t := range slotsTime {
			slots[i] = UTCTime(t)
		}
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

// возвращает список заказов опредеоённого киента
func (m *OrderManager) GetByClient(email string) ([]*ClientOrder, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}

	fullOrders, err := m.storage.GetByClient(email)

	if err != nil {
		return nil, err
	}

	var responses []*ClientOrder

	for _, fo := range fullOrders {
		priceMap := make(map[string]float32)
		for _, p := range fo.Price {
			priceMap[p.Detail] = p.Price
		}

		// Формируем ServDetails, объединяя длительность из OrderDetails и цену из Price
		servDetails := make([]ServDetails, 0, len(fo.OrderDetails))
		for _, d := range fo.OrderDetails {
			price, exists := priceMap[d.Detail]
			if !exists {
				// Если цены нет, деталь пропускаем (логически такого быть не должно)
				continue
			}
			servDetails = append(servDetails, ServDetails{
				Detail:   d.Detail,
				Duration: d.Duration,
				Price:    price,
			})
		}

		// Преобразование времени
		start := UTCTime(fo.StartMoment)
		var end *UTCTime
		if fo.EndMoment != nil {
			utcEnd := UTCTime(*fo.EndMoment)
			end = &utcEnd
		}

		responses = append(responses, &ClientOrder{
			ID:           fo.ID,
			NameCompany:  fo.NameCompany,
			City:         fo.City,
			Address:      fo.Address,
			Service:      fo.Service,
			StartMoment:  start,
			EndMoment:    end,
			Status:       fo.Status,
			OrderDetails: servDetails,
			Sum:          fo.Sum,
		})
	}

	return responses, nil
}

// func (m *OrderManager) GetByCompany(inn string) ([]*FullOrder, error) {
// 	if len(inn) != 12 {
// 		return nil, ErrInnLen
// 	}
// 	return m.storage.GetByCompany(inn)
// }

// // возвращает список всех заказов
// func (m *OrderManager) GetFullAllOrders() ([]*FullOrder, error) {
// 	return m.storage.GetFullAllOrders()
// }

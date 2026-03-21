package swagger

import (
	"src/internal/branch"
	"src/internal/client"
	"src/internal/order"
	"src/internal/service"
)

// GetClientOrders возвращает список заказов...
// @Summary      Получить заказы клиента
// @Description  Возвращает все заказы, принадлежащие авторизованному клиенту.
// @Tags         client
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  order.ClientOrderResponse
// @Failure      401  {string}  string
// @Failure      500  {string}  string  "Internal server error"
// @Router       /client/orders [get]
func getClientOrders() {
	// фиктивное использование, чтобы избежать ошибки "imported and not used"
	var _ = order.ClientOrderResponse{}
}

// GetServices - возвращает список сервисов
// @Summary Получить список услуг
// @Description Возвращает список общих услуг и их ID
// @Tags  services
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array} service.ServiceResponse
// @Failure      401  {string}  string
// @Failure      500  {string}  string  "Internal server error"
// @Router       /services [get]
func getServices() {
	// фиктивное использование, чтобы избежать ошибки "imported and not used"
	var _ = service.ServiceResponse{}
}

// GetClientCity
// @Summary Получить город клиента
// @Description Возвращает город авторизованного клиента
// @Tags client
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} client.GetCityResponse
// @Failure      401  {string}  string
// @Failure      500  {string}  string  "Internal server error"
// @Router       /client/city [get]
func getClientCity() {
	var _ = client.GetCityResponse{}
}

// PutClientCity
// @Summary Обновить город клиента
// @Description Обновляет город клиента, если тот есть в списке городов (в любом регистре)
// @Tags client
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} client.UpdateCityRequest
// @Failure      401  {string}  string
// @Failure      500  {string}  string  "Internal server error"
// @Router       /client/city [put]
func putClientCity() {
	var _ = client.UpdateCityRequest{}
}

// GetBranchesByCityAndService godoc
// @Summary      Получить филиал для заказа
// @Description  Получить филиал для определённого города с определённым сервисом
// @Tags         branch
// @Security     BearerAuth
// @Produce      json
// @Param        city    query string true  "Город"
// @Param        service query string true  "ID Сервиса"
// @Success      200  {array} branch.BrancByCityServ
// @Failure      401  {string} string
// @Failure      500  {string} string  "Internal server error"
// @Router       /branch [get]
func getBranchByCityServ() {
	var _ = branch.BrancByCityServ{}
}

// Get branch/service/details/{id_branchserv}
// @Summary      Получить состав услуги
// @Description  Получение состава услуги по id branchserv для авторизованных пользователей
// @Tags         branch
// @Security     BearerAuth
// @Produce      json
// @Param        id_branchserv  path string true  "branchserv ID"
// @Success      200  {array} branch.ServiceDetails
// @Failure      401  {string} string
// @Failure      500  {string} string  "Internal server error"
// @Router       /branch/service/details/{id_branchserv} [get]
func getServiceDetails() {
	var _ = branch.ServiceDetails{}
}

// Get /branch/freetime?branch_id=<id_branch>&date=<date>&duration=<сумма_минут_выбранных_услуг>
// @Summary      Получить доступное время для записи
// @Description  Возвращает доступное время записи в определёный филиал с учётом общей длительности услуги на неделю
// @Tags         branch
// @Security     BearerAuth
// @Produce      json
// @Param        branch_id    query string true  "ID Филиала"
// @Param        date query string true  "Дата начала недели формат yyyy-mm-dd"
// @Param        duration query string true  "Общая длительность услуги в минутах"
// @Success      200  {array} order.DailySlots
// @Failure      401  {string} string
// @Failure      500  {string} string  "Internal server error"
// @Router       /branch/freetime [get]
func getBranchFreeTime() {
	var _ = order.DailySlots{}
}

// CreateOrder создаёт новый заказ
// @Summary      Создать заказ
// @Description  Создаёт новый заказ для авторизованного клиента на доступное время.
// @Tags         order
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body order.CreateOrderRequest true "Данные для создания заказа"
// @Success      201  {object}  order.Order
// @Failure      400  {string}  string  "Invalid request body or business logic error"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /order [post]
func _createOrder() {
	var _ = order.CreateOrderRequest{}
	var _ = order.Order{}

}

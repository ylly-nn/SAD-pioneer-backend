package swagger

import "src/internal/admin"

// getAllPartnerRequests возвращает список всех заявок на регистрацию организаций
// @Summary      Получить все заявки партнёров
// @Description  Возвращает список всех заявок на регистрацию организаций. Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   admin.PartnerRequest  "Список заявок (если нет, возвращается null)"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      403  {string}  string  "Forbidden: admin access required"
// @Failure      500  {string}  string  "Failed to get requests"
// @Router       /admin/partner-requests/ [get]
func getAllPartnerRequests() {
	var _ = admin.PartnerRequest{}
}

// getNewPartnerRequests возвращает список новых заявок на регистрацию организаций
// @Summary      Получить новые заявки партнёров
// @Description  Возвращает список заявок со статусом "new". Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   admin.PartnerRequest  "Список новых заявок (если нет, возвращается null)"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      403  {string}  string  "Forbidden: admin access required"
// @Failure      500  {string}  string  "Failed to get requests"
// @Router       /admin/partner-requests/new [get]
func getNewPartnerRequests() {
	var _ = admin.PartnerRequest{}
}

// getPendingPartnerRequests возвращает список заявок в статусе "pending"
// @Summary      Получить заявки в работе
// @Description  Возвращает список заявок со статусом "pending" (в работе). Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   admin.PartnerRequest  "Список заявок в работе (если нет, возвращается null)"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      403  {string}  string  "Forbidden: admin access required"
// @Failure      500  {string}  string  "Failed to get requests"
// @Router       /admin/partner-requests/pending [get]
func getPendingPartnerRequests() {
	var _ = admin.PartnerRequest{}
}

// getApprovedPartnerRequests возвращает список принятых заявок
// @Summary      Получить принятые заявки партнёров
// @Description  Возвращает список заявок со статусом "approved" (принятые). Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   admin.PartnerRequest  "Список принятых заявок (если нет, возвращается null"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      403  {string}  string  "Forbidden: admin access required"
// @Failure      500  {string}  string  "Failed to get requests"
// @Router       /admin/partner-requests/approved [get]
func getApprovedPartnerRequests() {
	var _ = admin.PartnerRequest{}
}

// getRejectedPartnerRequests возвращает список отклонённых заявок
// @Summary      Получить отклонённые заявки партнёров
// @Description  Возвращает список заявок со статусом "rejected" (отклонённые). Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   admin.PartnerRequest  "Список отклонённых заявок (если нет, возвращается null)"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      403  {string}  string  "Forbidden: admin access required"
// @Failure      500  {string}  string  "Failed to get requests"
// @Router       /admin/partner-requests/rejected [get]
func getRejectedPartnerRequests() {
	var _ = admin.PartnerRequest{}
}

// takeRequestToWork переводит заявку в статус "в работе"
// @Summary      Взять заявку в работу
// @Description  Администратор переводит заявку из статуса "new" в "pending". Требуется указать INN организации.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      admin.ApprovePartnerRequest  true  "INN организации"
// @Success      200      {object}  map[string]string  "message: Request taken to work, inn: ..., status: pending"
// @Failure      400      {string}  string  "Invalid request body | validation error | request with inn ... not found | request cannot be taken to work: current status is ..."
// @Failure      401      {string}  string  "Unauthorized"
// @Failure      403      {string}  string  "Forbidden: admin access required"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /admin/partner-requests/take [post]
func takeRequestToWork() {
	var _ = admin.ApprovePartnerRequest{}
}

// approvePartnerRequest одобряет заявку партнёра
// @Summary      Одобрить заявку партнёра
// @Description  Администратор одобряет заявку (статус "pending" -> "approved"). Создаётся компания и связывается с пользователем.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      admin.ApprovePartnerRequest  true  "INN организации"
// @Success      200      {object}  map[string]string  "message: Partner request approved"
// @Failure      400      {string}  string  "Invalid request body | validation error | request with inn ... not found | request already processed | failed to create company | failed to create partner user record | failed to update request status"
// @Failure      401      {string}  string  "Unauthorized"
// @Failure      403      {string}  string  "Forbidden: admin access required"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /admin/partner-requests/approve [post]
func approvePartnerRequest() {
	var _ = admin.ApprovePartnerRequest{}
}

// rejectPartnerRequest отклоняет заявку партнёра
// @Summary      Отклонить заявку партнёра
// @Description  Администратор отклоняет заявку (статус "pending" -> "rejected"). Требуется указать INN организации.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      admin.ApprovePartnerRequest  true  "INN организации"
// @Success      200      {object}  map[string]string  "message: Request rejected, inn: ..., status: rejected"
// @Failure      400      {string}  string  "Invalid request body | validation error | request with inn ... not found | request cannot be rejected: current status is ..."
// @Failure      401      {string}  string  "Unauthorized"
// @Failure      403      {string}  string  "Forbidden: admin access required"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /admin/partner-requests/reject [post]
func rejectPartnerRequest() {
	var _ = admin.ApprovePartnerRequest{}
}

type CreateAdminRequest struct {
	Email   string `json:"email" example:"example@gmail.com" validate:"required,email"`
	Name    string `json:"name" example:"Иван" validate:"required"`
	Surname string `json:"surname" example:"Иванов" validate:"required"`
}

// CreateAdminResponse — структура ответа при успешном создании администратора
type CreateAdminResponse struct {
	Message string `json:"message" example:"Admin created successfully"`
	Email   string `json:"email" example:"example@gmail.com"`
}

// CreateAdmin создает нового администратора
// @Summary      Создать нового администратора
// @Description  Создает нового администратора. Доступно только для существующих администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateAdminRequest true "Данные для создания администратора"
// @Success      201 {object} CreateAdminResponse "Администратор успешно создан"
// @Failure      400 {string} string "Invalid request body | validation error | email already exists"
// @Failure      401 {string} string "Unauthorized | Invalid token"
// @Failure      403 {string} string "Forbidden: only admins can create new admins"
// @Failure      500 {string} string "Internal server error"
// @Router       /admin/create-admin [post]
func CreateAdmin() {
	var _ = CreateAdminRequest{}
	var _ = CreateAdminResponse{}
}

// GetById возвращает заявку на регистрацию организации по её ID
// @Summary      Получить заявку по ID
// @Description  Возвращает заявку на регистрацию организации по указанному UUID. Доступно только для администраторов.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "UUID заявки" format(uuid)
// @Success      200 {object} admin.PartnerRequest "Информация о заявке"
// @Failure      400 {string} string "ID is required | ID must be UUID"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "Forbidden: admin access required"
// @Failure      404 {string} string "request with id ... not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /admin/partner-requests/{id} [get]
func GetById() {
	var _ = admin.PartnerRequest{}
}

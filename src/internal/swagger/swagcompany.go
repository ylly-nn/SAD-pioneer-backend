package swagger

import (
	"src/internal/company"
)

// getCompany возвращает информацию о компании текущего пользователя
// @Summary      Получить компанию пользователя
// @Description  Возвращает данные компании для авторизованного пользователя (только для партнёров).
// @Tags         company
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  company.CompanyResponse  "Данные компании"
// @Failure      401  {string}  string  "Unauthorized: missing or invalid token"
// @Failure      403  {string}  string  "User does not have a company"
// @Failure      404  {string}  string  "Company not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /company [get]
func getCompany() {
	var _ = company.CompanyResponse{}
}

// getCompanyBranches возвращает список филиалов компании авторизованного пользователя
// @Summary      Получить филиалы компании
// @Description  Возвращает список филиалов компании, к которой привязан авторизованный пользователь (только для партнёров).
// @Tags         company
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   company.CompanyBranch  "Список филиалов (если нет филиалов, возвращается null)"
// @Failure      401  {string}  string                 "Unauthorized: missing or invalid token"
// @Failure      403  {string}  string                 "User does not have a company"
// @Failure      500  {string}  string                 "Internal server error"
// @Router       /company/branches [get]
func getCompanyBranches() {
	var _ = company.CompanyBranch{}
}

// getCompanyBranchesBranchID возвращает филиал компании с услугами
// @Summary      Получить филиал компании по ID
// @Description  Возвращает детальную информацию о филиале компании текущего пользователя, включая список услуг.
// @Tags         company
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branch_id  path      string  true  "UUID филиала"  example(9eebb3b9-5b35-4007-9d4f-2f4141786b45)
// @Success      200        {object}  company.CompanyBranchWithServ  "Филиал с услугами (если не найден, возвращается null)"
// @Failure      400        {string}  string  "missing branch id или invalid branch id format: must be UUID"
// @Failure      401        {string}  string  "unauthorized: missing user claims или email not found in token"
// @Failure      403        {string}  string  "User does not have a company или User does not have access to the branch"
// @Failure      500        {string}  string  "internal server error или failed to encode response"
// @Router       /company/branches/{branch_id} [get]
func getCompanyBranchesBranchID() {
	var _ = company.ServiceInBranch{}
	var _ = company.CompanyBranchWithServ{}
}

// getCompanyBranchServiceBranchserv_id возвращает детали услуги филиала
// @Summary      Получить детали услуги филиала
// @Description  Возвращает описание и длительность услуги, связанной с указанным идентификатором услуги филиала (branch_serv_id).
// @Tags         company
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        branchServID  path      string  true  "UUID записи услуги филиала"  example(6fdd2352-ffc4-4140-b54c-67657f841c1c)
// @Success      200  {object}  company.ServDetails  "Детали услуги (если данные отсутствуют, возвращается null)"
// @Failure      400  {string}  string  "invalid branch service ID или invalid service details format"
// @Failure      401  {string}  string  "unauthorized: missing user claims или email not found in token"
// @Failure      403  {string}  string  "user is not a partner или branch service not available"
// @Failure      500  {string}  string  "internal server error или failed to encode response"
// @Router       /company/branch/service/{branchServID} [get]
func getCompanyBranchServiceBranchserv_id() {
	var _ = company.ServDetails{}
}

// getCompanyOrders возвращает список заказов компании по филиалам
// @Summary      Получить заказы компании
// @Description  Возвращает список заказов компании, сгруппированных по филиалам. Доступно только для партнёров.
// @Tags         company
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   company.CompanyBranchOrderResponse  "Список филиалов с заказами (если данных нет, возвращается null)"
// @Failure      401  {string}  string  "unauthorized: missing user claims или email not found in token"
// @Failure      403  {string}  string  "user is not a partner"
// @Failure      500  {string}  string  "internal server error или failed to encode response"
// @Router       /company/orders [get]
func getCompanyOrders() {
	var _ = company.CompanyBranchOrderResponse{}
	var _ = company.CompanyOrder{}
	var _ = company.ServDetails{}
}

package company

import (
	"encoding/json"
	"errors"
	"net/http"
	"src/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы для компаний
type Handler struct {
	company *CompanyManager
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(company *CompanyManager) *Handler {
	return &Handler{company: company}
}

// GetCompanies обрабатывает GET /companies.
func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.company.GetAllCompanies()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]CompanyResponse, len(companies))
	for i, s := range companies {
		resp[i] = CompanyResponse{INN: s.INN, KPP: s.KPP, OGRN: s.OGRN, OrgName: s.OrgName, OrgShortName: s.OrgShortName}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// При успешном удалении возвращает статус 204 No Content.
// Если компания не найдена, возвращает 404 Not Found.
func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	innStr := chi.URLParam(r, "inn")
	if innStr == "" {
		http.Error(w, "INN is required", http.StatusBadRequest)
		return
	}

	if len(innStr) != 12 {
		http.Error(w, "INN must be 12 characters", http.StatusBadRequest)
		return
	}

	err := h.company.DeleteCompany(innStr)
	if err != nil {
		if errors.Is(err, ErrCompanyNotFound) {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCompany обрабатывает GET /company.
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	// Получение email из claims (из тела токена)
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	company, err := h.company.GetCompany(email)
	if err != nil {
		if errors.Is(err, ErrUserNotPartner) {
			http.Error(w, "User does not have a company", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrCompanyNotFound) {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := CompanyResponse{
		INN:          company.INN,
		KPP:          company.KPP,
		OGRN:         company.OGRN,
		OrgName:      company.OrgName,
		OrgShortName: company.OrgShortName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateCompany обрабатывает POST /companies.
func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	company, err := h.company.CreateCompany(req.Company)
	if err != nil {
		if errors.Is(err, ErrCompanyAlreadyExists) {
			http.Error(w, "Company with this INN already exists", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CompanyResponse{
		INN:          company.INN,
		KPP:          company.KPP,
		OGRN:         company.OGRN,
		OrgName:      company.OrgName,
		OrgShortName: company.OrgShortName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Обрабатывает Get /company/branches
func (h *Handler) GetBranchesByUser(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	branches, err := h.company.GetBranchesByEmail(email)
	if err != nil {
		// Если пользователь не является партнёром
		if errors.Is(err, ErrUserNotPartner) {
			http.Error(w, "User does not have a company", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrBranchesNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil) // тело ответа: null
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(branches); err != nil {

		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Обрабатывает Get /company/branches/{branch_id}
func (h *Handler) GetBrancesByIdUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "branch_id")
	if idParam == "" {
		http.Error(w, "missing branch id", http.StatusBadRequest)
		return
	}

	// Парсим UUID и проверяем формат
	branchID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid branch id format: must be UUID", http.StatusBadRequest)
		return
	}

	branch, err := h.company.GetBranchByIdEmail(branchID, email)
	if err != nil {
		if errors.Is(err, ErrUserNotPartner) {
			http.Error(w, "User does not have a company", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrBranchNotInCompany) {
			http.Error(w, "User does not have access to the branch", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrBranchNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil) // тело ответа: null
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(branch); err != nil {
		// Если сериализация не удалась – ошибка сервера
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// обрабатывает get /company/branch/service/{branch_serv_id}
func (h *Handler) GetServDetailsByBranchServId(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}
	// Извлекаем branchServID из пути
	branchServIDStr := chi.URLParam(r, "branchServID")
	branchServID, err := uuid.Parse(branchServIDStr)
	if err != nil {
		http.Error(w, "invalid branch service ID", http.StatusBadRequest)
		return
	}

	// Вызываем бизнес-логику
	details, err := h.company.GetServDetailsByBranchServId(branchServID, email)
	if err != nil {
		// Обработка известных ошибок
		switch {
		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "user is not a partner", http.StatusForbidden)
		case errors.Is(err, ErrBranchesNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchServNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchServNotAvailable):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchServIsNull):
			{
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(nil) // тело ответа: null
				return
			}
		case errors.Is(err, ErrServiceDetailsInvalid):
			// Ошибка формата JSON (может быть как 400, так и 500, выбираем 400)
			http.Error(w, "invalid service details format", http.StatusBadRequest)
		default:
			// Неизвестная ошибка (БД, сеть и т.д.)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Обрабатывает Post /company/users
func (h *Handler) AddNewUserToCompany(w http.ResponseWriter, r *http.Request) {
	// Получение claims
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Извлечение email из claims
	userEmail, ok := claims["email"].(string)
	if !ok || userEmail == "" {
		http.Error(w, "Invalid token: email not found", http.StatusUnauthorized)
		return
	}

	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	// Добавление пользователя
	err := h.company.AddUserToCompany(userEmail, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "User does not have a company", http.StatusForbidden)
		case errors.Is(err, ErrCompanyNotFound):
			http.Error(w, "Company not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to company successfully",
		"email":   req.Email,
	})
}

// Обрабатывает Post /company/branch
func (h *Handler) AddNewBranchToCompany(w http.ResponseWriter, r *http.Request) {
	// Получение claims
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Извлечение email из claims
	userEmail, ok := claims["email"].(string)
	if !ok || userEmail == "" {
		http.Error(w, "Invalid token: email not found", http.StatusUnauthorized)
		return
	}

	var req AddBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	// Добавление филиала
	err := h.company.AddBranchToCompany(userEmail, req.City, req.Address, req.OpenTime, req.CloseTime)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "User does not have a company", http.StatusForbidden)
		case errors.Is(err, ErrCompanyNotFound):
			http.Error(w, "Company not found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Branch added to company successfully",
		"city":       req.City,
		"address":    req.Address,
		"open_time":  req.OpenTime,
		"close_time": req.CloseTime,
	})
}

// GetCompanyOrders обрабатывает GET /company/orders
func (h *Handler) GetCompanyOrders(w http.ResponseWriter, r *http.Request) {
	// Получаем claims из контекста
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	orders, err := h.company.GetCompanyOrders(email)
	if err != nil {
		// Обрабатываем известные ошибки
		if errors.Is(err, ErrUserNotPartner) {
			http.Error(w, "user is not a partner", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrBranchNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil) // тело ответа: null
			return
		}
		if errors.Is(err, ErrOrderNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil) // тело ответа: null
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Успешный ответ – всегда массив, даже пустой
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateOrderStatus обрабатывает PUT /company/order/status
// Параметры query: orderID (UUID), status (approve|reject)
func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	orderIDStr := r.URL.Query().Get("orderID")
	if orderIDStr == "" {
		http.Error(w, "missing orderID parameter", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, "invalid orderID format: must be UUID", http.StatusBadRequest)
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		http.Error(w, "missing status parameter", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.company.UpdateOrderStatus(email, orderID, statusStr)
	if err != nil {

		switch {
		case errors.Is(err, ErrStatus):
			http.Error(w, ErrStatus.Error(), http.StatusBadRequest)

		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "user is not a partner", http.StatusForbidden)

		case errors.Is(err, ErrOrderNotAvailable):
			http.Error(w, ErrOrderNotAvailable.Error(), http.StatusForbidden)

		case errors.Is(err, ErrOrderNotFound):
			http.Error(w, ErrOrderNotAvailable.Error(), http.StatusForbidden)

		case errors.Is(err, ErrBranchNotFound):
			http.Error(w, ErrOrderNotAvailable.Error(), http.StatusForbidden)

		case errors.Is(err, ErrUpdateStatus):
			http.Error(w, ErrUpdateStatus.Error(), http.StatusBadRequest)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// AddServDetail обрабатывает POST /company/branch/service/detail
func (h *Handler) AddServDetail(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	var req AddServDetailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateStruct(req); err != nil {
		middleware.SendValidationError(w, err)
		return
	}

	details := ServDetails{
		Detail:   req.Detail,
		Duration: req.Duration,
	}

	createdDetail, err := h.company.AddServiceDetail(req.BranchServID, email, details)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "user is not a partner", http.StatusForbidden)
		case errors.Is(err, ErrBranchServNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchNotInCompany):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchServDetailAlreadyExists):
			http.Error(w, "detail for this service in the branch already exists", http.StatusConflict)
		case errors.Is(err, ErrInvalidDuration):
			http.Error(w, "invalid duration", http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdDetail); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteServDetail обрабатывает DELETE /company/branch/service/detail/{branchServID}
func (h *Handler) DeleteServDetail(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: missing user claims", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized: email not found in token", http.StatusUnauthorized)
		return
	}

	branchServIDStr := chi.URLParam(r, "branchServID")
	if branchServIDStr == "" {
		http.Error(w, "missing branchServID parameter", http.StatusBadRequest)
		return
	}

	branchServID, err := uuid.Parse(branchServIDStr)
	if err != nil {
		http.Error(w, "invalid branch service ID format: must be UUID", http.StatusBadRequest)
		return
	}

	detailName := r.URL.Query().Get("detail")
	if detailName == "" {
		http.Error(w, "missing detail parameter", http.StatusBadRequest)
		return
	}

	updatedDetails, err := h.company.DeleteServiceDetail(branchServID, email, detailName)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotPartner):
			http.Error(w, "user is not a partner", http.StatusForbidden)
		case errors.Is(err, ErrBranchServNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchNotFound):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrBranchNotInCompany):
			http.Error(w, "branch service not available", http.StatusForbidden)
		case errors.Is(err, ErrDetailNotFound):
			http.Error(w, "detail not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 5. Успешный ответ – обновлённый список деталей
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedDetails); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

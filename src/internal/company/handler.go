package company

import (
	"encoding/json"
	"errors"
	"net/http"

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

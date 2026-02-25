package company

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler обрабатывает HTTP-запросы для компаний
type Handler struct {
	company *CompanyManager
}

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
		// Внутренняя ошибка сервера
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCompanyByInn(w http.ResponseWriter, r *http.Request) {
	innStr := chi.URLParam(r, "inn")
	if innStr == "" {
		http.Error(w, "INN is required", http.StatusBadRequest)
		return
	}

	if len(innStr) != 12 {
		http.Error(w, "INN must be 12 characters", http.StatusBadRequest)
		return
	}

	company, err := h.company.GetCompanyByInn(innStr)
	if err != nil {
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
		// Ошибки валидации (обязательные поля)
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

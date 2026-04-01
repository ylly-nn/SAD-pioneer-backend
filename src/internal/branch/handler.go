package branch

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы, связанные с branch_services.
type Handler struct {
	branch *BranchManager
}

// NewHandler создаёт новый Handler с заданным BranchManager.
func NewHandler(branch *BranchManager) *Handler {
	return &Handler{branch: branch}
}

// CreateBranchService обрабатывает POST /company/branchserv.
func (h *Handler) CreateBranchService(w http.ResponseWriter, r *http.Request) {
	var req CreateBranchServRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	bs, err := h.branch.CreateBranchService(req)
	if err != nil {
		log.Printf("CreateBranchService error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bs); err != nil {
		log.Printf("CreateBranchService encode error: %v", err)
	}
}

// CreateBranch обрабатывает POST /branch.
func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var req CreateBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	branch, err := h.branch.CreateBranch(req)
	if err != nil {
		log.Printf("CreateBranch error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(branch); err != nil {
		log.Printf("CreateBranch encode error: %v", err)
	}
}

// GetBranchesByCityAndService обрабатывает GET /branch?city=...&service=...
func (h *Handler) GetBranchesByCityAndService(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из query
	city := r.URL.Query().Get("city")
	service := r.URL.Query().Get("service")

	if city == "" || service == "" {
		http.Error(w, "city and service parameters are required", http.StatusBadRequest)
		return
	}

	branches, err := h.branch.GetBranchByCityServ(city, service)
	if err != nil {
		log.Printf("GetBranchByCityServ error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если результат nil или пустой, можно вернуть пустой массив, а не null
	if branches == nil {
		branches = []*BrancByCityServ{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(branches); err != nil {
		log.Printf("GetBranchesByCityAndService encode error: %v", err)
	}
}

// GetServiceDetails обрабатывает GET /service/details/{id_branchserv}
func (h *Handler) GetServiceDetails(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id_branchserv")
	if idStr == "" {
		http.Error(w, "missing branch service id", http.StatusBadRequest)
		return
	}

	// Парсим UUID
	branchServID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid branch service id format", http.StatusBadRequest)
		return
	}

	details, err := h.branch.GetServiceDetails(branchServID)
	if err != nil {

		if errors.Is(err, ErrBranchServiceNotFound) {
			http.Error(w, "branch service not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrDetailsNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("GetServiceDetails encode error: %v", err)
	}
}

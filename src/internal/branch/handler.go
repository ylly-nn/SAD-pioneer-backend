package branch

import (
	"encoding/json"
	"log"
	"net/http"
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

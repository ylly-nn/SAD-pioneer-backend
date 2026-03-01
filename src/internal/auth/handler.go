package auth

import (
	"encoding/json"
	"net/http"

	"src/internal/middleware"
)

// Handler обрабатывает HTTP-запросы для авторизации пользователей
type Handler struct {
	auth *AuthManager
}

// Cоздаёт новый экземпляр Handler
func NewHandler(service *AuthManager) *Handler {
	return &Handler{auth: service}
}

// Register обрабатывает POST /auth/register, выполняет регистрацию
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Валидируем с помощью middleware
	if err := middleware.ValidateStruct(req); err != nil {
		// middleware сам отправит структурированный ответ с ошибками
		middleware.SendValidationError(w, err)
		return
	}

	err := h.auth.Register(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case ErrUserAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification code sent to email",
		"email":   req.Email,
	})
}

// VerifyCode обрабатывает POST /auth/verify, отправляет код подтверждения
func (h *Handler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Code == "" {
		http.Error(w, "Email and code are required", http.StatusBadRequest)
		return
	}

	err := h.auth.VerifyCode(req.Email, req.Code)
	if err != nil {
		switch err.Error() {
		case ErrInvalidCode, ErrCodeExpired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User successfully registered",
	})
}

// Login обрабатывает POST /auth/login, выполняет вход в аккаунт
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	tokens, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case ErrUserNotFound, ErrInvalidPassword:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

// Refresh обрабатывает POST /auth/refresh, работает с токенами
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	tokens, err := h.auth.RefreshTokens(req.RefreshToken)
	if err != nil {
		switch err.Error() {
		case ErrInvalidRefreshToken, ErrRefreshTokenExpired, ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

// Logout обрабатывает POST /auth/logout, выполняет выход из аккаунта
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получение refresh токена из тела запроса
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		// Пробуем получить из заголовка (для обратной совместимости)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			req.RefreshToken = authHeader[7:]
		} else {
			http.Error(w, "Refresh token is required", http.StatusBadRequest)
			return
		}
	}

	if err := h.auth.Logout(req.RefreshToken); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

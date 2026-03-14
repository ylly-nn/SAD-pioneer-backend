package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Интерфейс для проверки админа
type AdminChecker interface {
	IsAdmin(email string) (bool, error)
}

type AdminMiddleware struct {
	adminChecker AdminChecker
}

func NewAdminMiddleware(adminChecker AdminChecker) *AdminMiddleware {
	return &AdminMiddleware{adminChecker: adminChecker}
}

func (m *AdminMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		isAdmin, err := m.adminChecker.IsAdmin(email)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "is_admin", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

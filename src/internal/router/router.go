package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"src/internal/service"
)

func New(serviceHandler *service.Handler) http.Handler {
	r := chi.NewRouter()

	// Глобальные middleware для всех запросов
	r.Use(middleware.Logger)    // логирование запросов
	r.Use(middleware.Recoverer) // восстановление после паник

	// Маршруты для работы с услугами
	r.Route("/services", func(r chi.Router) {
		r.Get("/", serviceHandler.GetServices)          // GET /services — список
		r.Post("/", serviceHandler.CreateService)       // POST /services — создание
		r.Delete("/{id}", serviceHandler.DeleteService) // DELETE /services/{id} — удаление
	})

	return r
}

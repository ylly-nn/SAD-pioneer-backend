package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"src/internal/client"
	"src/internal/company"
	"src/internal/order"
	"src/internal/service"
)

func New(serviceHandler *service.Handler, companyHandler *company.Handler, clientHandler *client.Handler, orderHandler *order.Handler) http.Handler {
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

	r.Route("/company", func(r chi.Router) {
		r.Get("/{inn}", companyHandler.GetCompanyByInn)
		r.Get("/", companyHandler.GetCompanies)
		r.Post("/", companyHandler.CreateCompany)
		r.Delete("/{inn}", companyHandler.DeleteCompany)
	})

	r.Route("/client", func(r chi.Router) {
		r.Post("/", clientHandler.CreateClient)
		r.Put("/city", clientHandler.UpdateCity)
		r.Get("/city/{email}", clientHandler.GetCity)
	})

	r.Route("/order", func(r chi.Router) {
		r.Get("/", orderHandler.GetAllOrders)
	})

	return r
}

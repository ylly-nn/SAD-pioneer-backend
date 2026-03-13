package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"src/internal/auth"
	"src/internal/branch"
	"src/internal/client"
	"src/internal/company"
	"src/internal/middleware"
	"src/internal/order"
	"src/internal/service"
)

func New(authMiddleware *middleware.AuthMiddleware, serviceHandler *service.Handler, companyHandler *company.Handler, clientHandler *client.Handler, orderHandler *order.Handler, branchHandler *branch.Handler, authHandler *auth.Handler) http.Handler {
	r := chi.NewRouter()

	// Глобальные middleware для всех запросов
	r.Use(chimiddleware.Logger)    // логирование запросов
	r.Use(chimiddleware.Recoverer) // восстановление после паник

	// Маршруты для работы с услугами
	r.Route("/services", func(r chi.Router) {
		r.Post("/", serviceHandler.CreateService)
		r.Delete("/{id}", serviceHandler.DeleteService)

		//Защищённые маршруты
		r.With(authMiddleware.Authenticate).Get("/", serviceHandler.GetServices)
	})

	r.Route("/company", func(r chi.Router) {
		r.Get("/{inn}", companyHandler.GetCompanyByInn)
		r.Get("/", companyHandler.GetCompanies)
		r.Post("/", companyHandler.CreateCompany)
		r.Delete("/{inn}", companyHandler.DeleteCompany)
		r.Get("/order/{inn}", orderHandler.GetCompanyOrders)
		r.Post("/branch/service", branchHandler.CreateBranchService)
	})

	r.Route("/client", func(r chi.Router) {
		r.Post("/", clientHandler.CreateClient)

		//Защищённые маршруты
		r.With(authMiddleware.Authenticate).Put("/city", clientHandler.UpdateCity)
		r.With(authMiddleware.Authenticate).Get("/city", clientHandler.GetCity)
		r.With(authMiddleware.Authenticate).Get("/orders", orderHandler.GetClientOrders)
	})

	r.Route("/order", func(r chi.Router) {
		r.Get("/", orderHandler.GetFullAllOrders)

		//Защищённые маршруты
		r.With(authMiddleware.Authenticate).Post("/", orderHandler.CreateOrder)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/verify", authHandler.VerifyCode)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
	})

	r.Route("/branch", func(r chi.Router) {
		r.Post("/", branchHandler.CreateBranch)

		//Защищённые маршруты
		r.With(authMiddleware.Authenticate).Get("/freetime", orderHandler.GetFreeTime)
		r.With(authMiddleware.Authenticate).Get("/service/details/{id_branchserv}", branchHandler.GetServiceDetails)
		r.With(authMiddleware.Authenticate).Get("/", branchHandler.GetBranchesByCityAndService)
	})

	return r
}

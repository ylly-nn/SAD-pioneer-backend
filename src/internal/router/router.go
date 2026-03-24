package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"src/internal/admin"
	"src/internal/auth"
	"src/internal/branch"
	"src/internal/client"
	"src/internal/company"
	"src/internal/middleware"
	"src/internal/order"
	"src/internal/partners"
	"src/internal/service"
)

func New(authMiddleware *middleware.AuthMiddleware, adminMiddleware *middleware.AdminMiddleware, serviceHandler *service.Handler, companyHandler *company.Handler, clientHandler *client.Handler, orderHandler *order.Handler, branchHandler *branch.Handler, authHandler *auth.Handler, adminHandler *admin.Handler, partnersHandler *partners.Handler) http.Handler {
	r := chi.NewRouter()

	// Глобальные middleware для всех запросов
	r.Use(chimiddleware.Logger)    // логирование запросов
	r.Use(chimiddleware.Recoverer) // восстановление после паник

	// Swagger UI — доступен по /swagger/index.html
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.InstanceName("swagger"),
	))

	// Маршруты для работы с услугами
	r.Route("/services", func(r chi.Router) {
		r.Post("/", serviceHandler.CreateService)
		r.Delete("/{id}", serviceHandler.DeleteService)

		//Защищённые маршруты
		r.With(authMiddleware.Authenticate).Get("/", serviceHandler.GetServices)
	})

	r.Route("/company", func(r chi.Router) {
		r.Post("/", companyHandler.CreateCompany)
		r.Delete("/{inn}", companyHandler.DeleteCompany)
		r.Get("/order/{inn}", orderHandler.GetCompanyOrders)
		r.Post("/branch/service", branchHandler.CreateBranchService)

		//Защищёные маршруты
		r.With(authMiddleware.Authenticate).Get("/", companyHandler.GetCompany)
		r.With(authMiddleware.Authenticate).Get("/branches", companyHandler.GetBranchesByUser)
		r.With(authMiddleware.Authenticate).Get("/branches/{branch_id}", companyHandler.GetBrancesByIdUser)
		r.With(authMiddleware.Authenticate).Get("/branch/service/{branchServID}", companyHandler.GetServDetailsByBranchServId)
		r.With(authMiddleware.Authenticate).Post("/users", companyHandler.AddNewUserToCompany)
		r.With(authMiddleware.Authenticate).Post("/branch", companyHandler.AddNewBranchToCompany)
		r.With(authMiddleware.Authenticate).Get("/orders", companyHandler.GetCompanyOrders)
		r.With(authMiddleware.Authenticate).Put("/order/status", companyHandler.UpdateOrderStatus)
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

	r.Route("/partner", func(r chi.Router) {
		// Защищенные маршруты
		r.Use(authMiddleware.Authenticate)
		r.Post("/request", partnersHandler.CreatePartnerRequest)
		r.Get("/request", partnersHandler.GetRequestStatus)
	})

	r.Route("/admin", func(r chi.Router) {
		// Защищенные маршруты + проверка на админа
		r.Use(authMiddleware.Authenticate)
		r.Use(adminMiddleware.RequireAdmin)

		r.Get("/partner-requests/", adminHandler.GetAllRequests)
		r.Get("/partner-requests/new", adminHandler.GetNewRequests)
		r.Get("/partner-requests/pending", adminHandler.GetPendingRequests)
		r.Get("/partner-requests/approved", adminHandler.GetApprovedRequests)
		r.Get("/partner-requests/rejected", adminHandler.GetRejectedRequests)
		r.Post("/partner-requests/take", adminHandler.TakeRequestToWork)
		r.Post("/partner-requests/approve", adminHandler.ApprovePartnerRequest)
		r.Post("/partner-requests/reject", adminHandler.RejectPartnerRequest)
	})

	return r
}

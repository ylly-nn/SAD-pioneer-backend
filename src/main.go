package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"src/internal/auth"
	"src/internal/branch"
	"src/internal/client"
	"src/internal/company"
	configPkg "src/internal/config"
	"src/internal/db"
	"src/internal/order"
	"src/internal/router"
	"src/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment")
	}

	// Загрузка порта из env
	port := os.Getenv("SERVER_PORT")

	// Загрузка конфигурации для jwt токенов
	jwt, err := configPkg.LoadJWTConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	//Подключение к бд
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to database: %w", err)

	}
	defer database.Close()

	// Подключение к сервису с email
	emailService, err := configPkg.ConnectSMTP()
	if err != nil {
		log.Fatal("Failed to connect to SMTP:", err)
	}

	// Конфигурация токенов
	authConfig := auth.Config{
		JWTSecretKey:    jwt.SecretKey,
		AccessTokenTTL:  jwt.AccessTokenTTL,
		RefreshTokenTTL: jwt.RefreshTokenTTL,
		VerificationTTL: jwt.VerificationTTL,
	}

	//Запуск обработчиков из пакета servise
	serviceStorage := service.NewPostgresServiceStorage(database)
	serviceManager := service.NewServiceManager(serviceStorage)
	serviceHandler := service.NewHandler(serviceManager)

	//Запуск обработчиков из пакета company
	companyStorage := company.NewPostgresCompanyStorage(database)
	companyManager := company.NewCompanyManager(companyStorage)
	companyHandler := company.NewHandler(companyManager)

	clientStorage := client.NewPostgresClientStorage(database)
	clientManager := client.NewClientManager(clientStorage)
	clientHandler := client.NewHandler(clientManager)

	orderStorage := order.NewPostrgesOrderStorage(database)
	orderManager := order.NewOrderManager(orderStorage)
	orderHandler := order.NewHandler(orderManager)

	branchStorage := branch.NewPostgresBranchStorage(database)
	branchManager := branch.NewBranchManager(branchStorage)
	branchHandler := branch.NewHandler(branchManager)

	// Запуск обработчиков из пакета auth
	userStorage := auth.NewPostgresUserStorage(database)
	refreshTokenStorage := auth.NewPostgresRefreshTokenStorage(database)
	verificationStorage := auth.NewMemoryVerificationStorage()
	authService := auth.NewAuthManager(userStorage, refreshTokenStorage, verificationStorage, emailService, authConfig)
	authHandler := auth.NewHandler(authService)

	//Пути - src/internal/router/router.go
	router := router.New(serviceHandler, companyHandler, clientHandler, orderHandler, branchHandler, authHandler)

	//FIX(yumi): порт вынесен в .env (SERVER_PORT = 8080 (по умолчанию))
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	//Закрытие бд при выходе из main

}

package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"src/internal/client"
	"src/internal/company"
	"src/internal/db"
	"src/internal/order"
	"src/internal/router"
	"src/internal/service"
)

func main() {

	//Подключение к бд
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to database: %w", err)

	}
	defer database.Close()

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

	//Пути - src/internal/router/router.go
	router := router.New(serviceHandler, companyHandler, clientHandler, orderHandler)

	//TODO(ylly): вынести в .env Port
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	//Закрытие бд при выходе из main

}

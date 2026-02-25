package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"src/internal/db"
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

	storage := service.NewPostgresServiceStorage(database)
	serviceManager := service.NewServiceManager(storage)
	serviceHandler := service.NewHandler(serviceManager)
	router := router.New(serviceHandler)

	//TODO(ylly): вынести в .env Port
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	//Закрытие бд при выходе из main

}

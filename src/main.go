package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"src/internal/db"
)

func main() {
	r := mux.NewRouter()

	// FIX(ylly): потом удалить Тестовый эндпоинт
	r.HandleFunc("/test", testHandler).Methods("GET")

	//Подключение к бд
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to database: %w", err)

	}

	//TODO(ylly): вынести в .env Port
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	//Закрытие бд при выходе из main
	defer database.Close()
}

// FIX(ylly): потом удалить
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Storage struct {
	DB *sql.DB
}

// Содание структуры которая испоотзуется в других storage
func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

// Connect устанавливает соединение с БД и возвращает объект *sql.DB.
func Connect() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Формируем строку подключения
	// Для pgx/stdlib используется формат PostgreSQL-совместимый
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)

	// Открываем соединение (фактически создаётся пул соединений)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем, что соединение действительно работает
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return db, nil
}

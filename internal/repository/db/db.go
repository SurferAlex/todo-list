package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

var (
	// Переменная для хранения полдключения к базе данных
	db *sql.DB
)

// InitDB инициализирует подключение к базе данных
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем подключение к базе данных
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Успешно подключено к базе данных:", dataSourceName)
}

func GetDB() *sql.DB {
	return db
}

package main

import (
	"fmt"
	"os"
	"testi/internal/cmd"
	"testi/internal/repository/db"
	"testi/internal/session"
	"time"
)

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort,
	)

	// Подключение к БД
	db.InitDB(connectionString)

	// Создание таблиц в БД
	db.CreateUsersTable()
	db.CreateTasksTable()
	db.CreatePostsTable()
	db.CreatePostsImagesTable()

	// Инициализация Redis
	session.InitRedis()

	// Запуск сервера
	server := cmd.NewServer()
	server.Start()

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

// .

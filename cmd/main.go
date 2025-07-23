package main

import (
	"fmt"
	"testi/internal/cmd"
	"testi/internal/repository/db"
	"time"
)

func main() {
	// Инициализация подключения к базе данных
	db.InitDB(
		`user=postgres password=qwe1144EodT5
		dbname=test_blog
		host=localhost
		port=5433
		sslmode=disable`,
	)

	server := cmd.NewServer()
	server.Start()

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

// .

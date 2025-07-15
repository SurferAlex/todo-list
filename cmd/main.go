package main

import (
	"fmt"

	"testi/internal/cmd"
	"testi/internal/repository/db"
	"testi/internal/repository/tasks"

	"time"
)

func main() {
	// Инициализация подключения к базе данных
	db.InitDB("user=postgres password=Luc1an1998& dbname=user_db host=localhost port=5432 sslmode=disable") // Замените на ваши параметры подключения

	repo := tasks.NewRepoTasks() // путь к tasks.json встроен
	server := cmd.NewServer(repo)
	server.Start()

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

// .

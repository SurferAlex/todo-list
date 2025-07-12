package main

import (
	"fmt"

	"testi/internal/cmd"
	"testi/internal/repository/tasks"

	"time"
)

func main() {

	repo := tasks.NewRepoTasks() // путь к tasks.json встроен
	server := cmd.NewServer(repo)
	server.Start()

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

// .

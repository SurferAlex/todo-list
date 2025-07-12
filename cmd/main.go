package main

import (
	"fmt"
	"testi/internal/cmd"
	"time"
)

func main() {
	go cmd.StartServer() // Запускаем сервер в горутине

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

// .

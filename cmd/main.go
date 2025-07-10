package main

import (
	"fmt"
	"time"
)

// Объявление переменной tasks

func main() {
	go StartServer() // Запускаем сервер в горутине

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(1 * time.Second)
	}
}

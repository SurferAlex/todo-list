package main

import (
	"fmt"
	"testi/internal/cmd" // Добавляем импорт пакета cmd
	"time"
)

// Объявление переменной tasks

func main() {
	go cmd.StartServer() // Запускаем сервер в горутине

	for {
		fmt.Println("Сервер запущен. Ожидание запросов...")
		time.Sleep(10 * time.Second)
	}
}

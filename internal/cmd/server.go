package cmd

import (
	"fmt"
	"net/http"
	"testi/internal/api/router"
)

type Server struct{}

// NewServer теперь не принимает аргументов
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	// Настройка маршрутов
	router.SetupRouters()

	// Старт сервера
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", &router.Router{})
}

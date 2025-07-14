package cmd

import (
	"fmt"
	"net/http"
	"testi/internal/api/router"
	"testi/internal/repository/tasks"
)

type Server struct {
	repo *tasks.Repository
}

func NewServer(repo *tasks.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) Start() {

	// Настройка маршрутов
	router.SetupRouters(s.repo)

	// Старт сервера
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", &router.Router{})
}

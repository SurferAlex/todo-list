package cmd

import (
	"fmt"
	"html/template"
	"net/http"

	"testi/internal/api/router"
	"testi/internal/entity"
	"testi/internal/repository/tasks"
	"testi/internal/usecases/auth"
)

type Server struct {
	repo *tasks.Repository
}

func NewServer(repo *tasks.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) Start() {
	// Обслуживание статики
	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))

	// Роуты
	router.Handle("GET", "/", s.handleRoot)
	router.Handle("GET", "/home", s.handleHome)
	router.Handle("GET", "/tasks", s.handleGetTasks)
	router.Handle("POST", "/tasks", s.handlePostTasks)
	router.Handle("GET", "/login", auth.LoginHandler)
	router.Handle("POST", "/login", auth.CheckPassword)
	router.Handle("GET", "/register", auth.RegisterHandler)
	router.Handle("POST", "/register", auth.Register)
	router.Handle("GET", "/logout", auth.LogoutHandler)

	// Старт сервера
	http.HandleFunc("/", router.ServeHTTP)
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/home.html"))
	tmpl.Execute(w, nil)
}

func (s *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	username, err := s.getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tasksList, err := s.repo.GetTasksByUser(username)
	if err != nil {
		http.Error(w, "Ошибка загрузки задач", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("frontend/templates/tasks.html"))
	tmpl.Execute(w, tasksList)
}

func (s *Server) handlePostTasks(w http.ResponseWriter, r *http.Request) {
	username, err := s.getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := entity.User{Username: username}
	r.ParseForm()

	switch {
	case r.FormValue("title") != "":
		// Добавление новой задачи
		if err := s.repo.CreateTask(user, r.FormValue("title")); err != nil {
			http.Error(w, "Не удалось сохранить задачу", http.StatusInternalServerError)
			return
		}
	case r.FormValue("toggleId") != "":
		// Переключение статуса задачи
		if err := s.repo.ToggleCompleteByID(username, r.FormValue("toggleId")); err != nil {
			http.Error(w, "Ошибка переключения статуса", http.StatusInternalServerError)
			return
		}
	case r.FormValue("deleteId") != "":
		// Удаление задачи
		if err := s.repo.DeleteTaskByID(username, r.FormValue("deleteId")); err != nil {
			http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func (s *Server) getUsernameFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		return "", fmt.Errorf("cookie not found")
	}
	return cookie.Value, nil
}

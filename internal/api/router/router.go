package router

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"testi/internal/entity"
	"testi/internal/repository/tasks"
	"testi/internal/usecases/auth"
)

type routeKey struct {
	Method string
	Path   string
}

var routes = make(map[routeKey]http.HandlerFunc)

// Зарегистрировать маршрут
func Handle(method string, path string, handler http.HandlerFunc) {
	routes[routeKey{Method: strings.ToUpper(method), Path: path}] = handler
}

// Роутер как основной обработчик
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := routes[routeKey{Method: r.Method, Path: r.URL.Path}]; ok {
		handler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Router структура, реализующая интерфейс http.Handler
type Router struct{}

// ServeHTTP реализует интерфейс http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Обработка статики
	if strings.HasPrefix(req.URL.Path, "/frontend/") {
		http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))).ServeHTTP(w, req)
		return
	}
	// Ваши маршруты
	if handler, ok := routes[routeKey{Method: req.Method, Path: req.URL.Path}]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func SetupRouters(repo *tasks.Repository) {
	Handle("GET", "/", handleRoot)
	Handle("GET", "/home", handleHome)
	Handle("GET", "/tasks", handleGetTasks(repo))
	Handle("POST", "/tasks", handlePostTasks(repo))
	Handle("GET", "/login", auth.LoginHandler)
	Handle("POST", "/login", auth.CheckPassword)
	Handle("GET", "/register", auth.RegisterHandler)
	Handle("POST", "/register", auth.Register)
	Handle("GET", "/logout", auth.LogoutHandler)

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/home.html"))
	tmpl.Execute(w, nil)
}

func handleGetTasks(repo *tasks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getUsernameFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tasksList, err := repo.GetTasksByUser(username)
		if err != nil {
			http.Error(w, "Ошибка загрузки задач", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("frontend/templates/tasks.html"))
		tmpl.Execute(w, tasksList)
	}
}

func handlePostTasks(repo *tasks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getUsernameFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user := entity.User{Username: username}
		r.ParseForm()

		switch {
		case r.FormValue("title") != "":
			if err := repo.CreateTask(user, r.FormValue("title")); err != nil {
				http.Error(w, "Не удалось сохранить задачу", http.StatusInternalServerError)
				return
			}
		case r.FormValue("toggleId") != "":
			if err := repo.ToggleCompleteByID(username, r.FormValue("toggleId")); err != nil {
				http.Error(w, "Ошибка переключения статуса", http.StatusInternalServerError)
				return
			}
		case r.FormValue("deleteId") != "":
			if err := repo.DeleteTaskByID(username, r.FormValue("deleteId")); err != nil {
				http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	}
}

func getUsernameFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		return "", fmt.Errorf("cookie not found")
	}
	return cookie.Value, nil
}

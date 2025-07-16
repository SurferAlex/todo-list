package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"testi/internal/entity"
	"testi/internal/repository/db"
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

func SetupRouters() {
	Handle("GET", "/", handleRoot)
	Handle("GET", "/home", handleHome)
	Handle("GET", "/tasks", handleGetTasks)
	Handle("POST", "/tasks", handlePostTasks)
	Handle("GET", "/login", auth.LoginHandler)
	Handle("POST", "/login", auth.CheckPassword)
	Handle("GET", "/register", auth.RegisterHandler)
	Handle("POST", "/register", auth.Register)
	Handle("GET", "/logout", auth.LogoutHandler)
	Handle("POST", "/delete_account", auth.DeleteAccountHandler)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/home.html"))
	tmpl.Execute(w, nil)
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tasksList, err := db.GetTasksByUser(username)
	if err != nil {
		http.Error(w, "Ошибка загрузки задач", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("frontend/templates/tasks.html"))
	tmpl.Execute(w, tasksList)
}

func handlePostTasks(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := entity.User{Username: username}
	r.ParseForm()

	switch {
	case r.FormValue("title") != "":
		task := entity.Task{Username: user.Username, Title: r.FormValue("title")}
		if err := db.InsertTask(task); err != nil {
			http.Error(w, "Не удалось сохранить задачу", http.StatusInternalServerError)
			return
		}
	case r.FormValue("toggleId") != "":
		id, err := strconv.Atoi(r.FormValue("toggleId"))
		if err == nil {
			db.ToggleCompleteByID(username, id)
		}
	case r.FormValue("deleteId") != "":
		id, err := strconv.Atoi(r.FormValue("deleteId"))
		if err == nil {
			db.DeleteTaskByID(username, id)
		}
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

func getUsernameFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		return "", fmt.Errorf("cookie not found")
	}
	return cookie.Value, nil
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Например, имя пользователя берём из query-параметра (?username=...)
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	tasks, err := db.GetTasksByUser(username)
	if err != nil {
		http.Error(w, "failed to get tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if task.Username == "" || task.Title == "" {
		http.Error(w, "username and title are required", http.StatusBadRequest)
		return
	}

	err := db.InsertTask(task)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

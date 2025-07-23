package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testi/internal/entity"
	"testi/internal/repository/db"
	"testi/internal/usecases/auth"
	"time"
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
	Handle("GET", "/register", auth.RegisterHandler)
	Handle("POST", "/register", auth.Register)
	Handle("GET", "/login", auth.LoginHandler)
	Handle("POST", "/login", auth.CheckPassword)
	Handle("GET", "/logout", auth.LogoutHandler)
	Handle("POST", "/delete_account", auth.DeleteAccountHandler)
	Handle("GET", "/tasks", handleGetTasks)
	Handle("POST", "/tasks", handlePostTasks)
	Handle("GET", "/wall", handleGetWall)
	Handle("POST", "/wall", handleAddPosts)
	Handle("POST", "/delete_post", handleDeletePost)
	Handle("GET", "/profile", handleProfile) // добавляем новый маршрут
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

func handleGetWall(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetAllPosts()
	if err != nil {
		fmt.Printf("Ошибка загрузки постов: %v\n", err) // добавь для отладки
		http.Error(w, "Ошибка загрузки постов", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("frontend/templates/wall.html"))
	tmpl.Execute(w, posts)
}

func handleAddPosts(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		fmt.Printf("Ошибка получения userID: %v\n", err)
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Парсим multipart форму
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("Ошибка парсинга формы: %v\n", err)
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	post := entity.Post{UserID: userID, Content: content}

	// Сначала создаём пост
	if err := db.InsertPost(post); err != nil {
		fmt.Printf("Ошибка создания поста: %v\n", err)
		http.Error(w, "Не удалось создать пост", http.StatusInternalServerError)
		return
	}

	// Получаем ID созданного поста
	postID, err := db.GetLastPostID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения ID поста: %v\n", err)
		http.Error(w, "Ошибка получения ID поста", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Создан пост с ID: %d\n", postID)

	// Обрабатываем загруженные файлы
	files := r.MultipartForm.File["images"]
	fmt.Printf("Загружено файлов: %d\n", len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("Ошибка открытия файла: %v\n", err)
			continue
		}
		defer file.Close()

		// Генерируем уникальное имя файла
		filename := generateUniqueFilename(fileHeader.Filename)
		filepath := "frontend/uploads/" + filename

		// Создаём папку если её нет
		os.MkdirAll("frontend/uploads", 0755)

		// Сохраняем файл
		dst, err := os.Create(filepath)
		if err != nil {
			fmt.Printf("Ошибка создания файла: %v\n", err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			fmt.Printf("Ошибка копирования файла: %v\n", err)
			continue
		}

		// Сохраняем информацию о файле в БД
		image := entity.PostImage{
			PostID:       postID,
			Filename:     filename,
			OriginalName: fileHeader.Filename,
			FilePath:     filepath,
		}
		if err := db.InsertPostImage(image); err != nil {
			fmt.Printf("Ошибка сохранения изображения в БД: %v\n", err)
		}
	}

	http.Redirect(w, r, "/wall", http.StatusSeeOther)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("frontend/templates/profile.html"))
	tmpl.Execute(w, username)
}

func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := db.DeletePostByID(postID, userID); err != nil {
		http.Error(w, "Не удалось удалить пост", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wall", http.StatusSeeOther)
}

func generateUniqueFilename(originalName string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%d%s", timestamp, ext)
}

package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"testi/internal/entity"
	"testi/internal/repository/db"
	"testi/internal/session"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка на существование пользователя
	existingUser, err := db.GetUserByUsername(creds.Username) // Функция для получения пользователя по имени
	if err == nil && existingUser != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	// Добавление нового пользователя в базу данных
	newUser := entity.User{Username: creds.Username, Password: creds.Password}
	if err := db.InsertUser(newUser); err != nil {
		http.Error(w, "Ошибка при сохранении пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/register.html"))
	tmpl.Execute(w, nil)
}

func CheckPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== Начало CheckPassword ===")

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		fmt.Printf("Ошибка декодирования: %v\n", err)
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	fmt.Printf("Попытка входа для пользователя: %s\n", creds.Username)

	// Проверка логина и пароля
	user, err := db.GetUserByUsername(creds.Username)
	if err != nil || user == nil || user.Password != creds.Password {
		fmt.Printf("Неверные данные для пользователя: %s\n", creds.Username)
		http.Error(w, "Неверные данные", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Пользователь найден: %s\n", creds.Username)

	// Получить userID
	userID, err := db.GetUserIDByUsername(creds.Username)
	if err != nil {
		fmt.Printf("Ошибка получения userID: %v\n", err)
		http.Error(w, "Ошибка получения данных пользователя", http.StatusInternalServerError)
		return
	}

	fmt.Printf("UserID получен: %d\n", userID)

	// Создать сессию
	sessionID, err := session.CreateSession(userID, creds.Username)
	if err != nil {
		fmt.Printf("Ошибка создания сессии: %v\n", err)
		http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Сессия создана: %s\n", sessionID)

	// Установить cookie с сессией
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 7, // 7 дней
		Secure:   false,            // для HTTP
		SameSite: http.SameSiteLaxMode,
	})

	fmt.Println("=== Успешный вход ===")
	// Успешный вход
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Отображение формы входа
	tmpl := template.Must(template.ParseFiles("frontend/templates/login.html"))
	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Получить session_id из cookie
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		// Удалить сессию из Redis
		session.DeleteSession(cookie.Value)
	}

	// Удалить cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Перенаправление на логин
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	err = db.DeleteUser(username) // Удаляем пользователя из базы данных
	if err != nil {
		http.Error(w, "Ошибка при удалении аккаунта", http.StatusInternalServerError)
		return
	}

	// Удаляем куку и редиректим на главную
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

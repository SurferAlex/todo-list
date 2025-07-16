package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
	"testi/internal/entity"
	"testi/internal/repository/db" //
	"time"
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
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка логина и пароля
	user, err := db.GetUserByUsername(creds.Username) // Получаем пользователя из базы данных
	if err != nil || user == nil || user.Password != creds.Password {
		http.Error(w, "Неверные данные", http.StatusUnauthorized)
		return
	}

	// Успешный вход
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Отображение формы входа
	tmpl := template.Must(template.ParseFiles("frontend/templates/login.html"))
	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Удаление куки username
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})

	// Перенаправление на логин или главную
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

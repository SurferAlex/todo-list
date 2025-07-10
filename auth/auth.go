package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type User struct {
	Username string
	Password string
}

var users []User

const usersFile = "users.json"

func loadUsers() error {
	file, err := os.Open(usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&users)
}

func saveUsers() error {
	file, err := os.Create(usersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(users)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		for _, user := range users {
			if user.Username == username {
				http.Error(w, "Пользователь уже существует", http.StatusConflict)
				return
			}
		}

		users = append(users, User{Username: username, Password: password})
		saveUsers()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("frontend/templates/register.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона регистрации", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		for _, user := range users {
			if user.Username == username && user.Password == password {
				// Успешный вход
				http.Redirect(w, r, "/tasks", http.StatusSeeOther)
				return
			}
		}

		http.Error(w, "Неверные данные", http.StatusUnauthorized)
		return
	}
	// Отображение формы входа
	tmpl := template.Must(template.ParseFiles("frontend/templates/login.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона входа", http.StatusInternalServerError)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь можно добавить логику для завершения сессии, если она есть
	// Например, удалить куки или сессии, если вы их используете

	// Перенаправление на главную страницу или страницу входа
	http.Redirect(w, r, "/glav", http.StatusSeeOther)
}

func init() {
	loadUsers()
}

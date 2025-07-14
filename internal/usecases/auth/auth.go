package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"testi/internal/entity"
)

var users []entity.User

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
		fmt.Println("Ошибка при создании файла:", err)
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(users); err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return err
	}
	fmt.Println("Пользователи успешно сохранены в файл.")
	return nil
}
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
	for _, user := range users {
		if user.Username == creds.Username {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
	}

	// Добавление нового пользователя
	users = append(users, entity.User{Username: creds.Username, Password: creds.Password})

	// Сохранение пользователей
	if err := saveUsers(); err != nil {
		fmt.Println("Ошибка при сохранении пользователей:", err)
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
	for _, user := range users {
		if user.Username == creds.Username && user.Password == creds.Password {
			// Успешный вход
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Неверные данные", http.StatusUnauthorized)
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

func DeleteUser(username string) error {
	err := loadUsers()
	if err != nil {
		return err
	}
	newUsers := make([]entity.User, 0, len(users))
	for _, user := range users {
		if user.Username != username {
			newUsers = append(newUsers, user)
		}
	}
	users = newUsers
	return saveUsers()
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	err = DeleteUser(username)
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

func init() {
	loadUsers()
}

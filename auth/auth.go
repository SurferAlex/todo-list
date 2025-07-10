package auth

import (
	"encoding/json"
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

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Регистрация</title>
	</head>
	<body>
		<h1>Регистрация</h1>
		<form method="post">
			<input type="text" name="username" placeholder="Имя пользователя" required>
			<input type="password" name="password" placeholder="Пароль" required>
			<button type="submit">Зарегистрироваться</button>
		</form>
	</body>
	</html>`
	w.Write([]byte(tmpl))
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
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>Вход</title>
</head>
<body>
	<h1>Вход</h1>
	<form method="post">
		<input type="text" name="username" placeholder="Имя пользователя" required>
		<input type="password" name="password" placeholder="Пароль" required>
		<button type="submit">Войти</button>
	</form>
</body>
</html>`
	w.Write([]byte(tmpl))
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

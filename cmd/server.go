package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"testi/auth/auth"
	"time" // Добавляем пакет time для работы с датами
)

type Task struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Deadline  string `json:"deadline"`
	Overdue   bool   `json:"overdue"` // Новое поле для статуса просрочки
}

var tasks []Task

// Функция для сохранения задач в файл
func saveTasksToFile(filename string) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при сохранении файла:", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
	}
}

// Функция для загрузки задач из файла
func loadTasksFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return // файла нет — просто начинаем с пустого списка
		}
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("Ошибка при разборе задач:", err)
		return
	}

	// Устанавливаем статус просрочки
	for i := range tasks {
		if tasks[i].Deadline != "" {
			deadline, err := time.Parse("2006-01-02", tasks[i].Deadline)
			if err == nil && time.Now().After(deadline) {
				tasks[i].Overdue = true
			}
		}
	}
}

// Функция для форматирования даты
func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr // Возвращаем оригинальную строку в случае ошибки
	}
	return date.Format("02-01-2006") // Форматируем дату как DD-MM-YYYY
}

// Обработчик для отображения новых задач и добавления новых задач
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		if title := r.FormValue("title"); title != "" {
			// Обработка добавления новой задачи
			deadline := r.FormValue("deadline")
			tasks = append(tasks, Task{Title: title, Completed: false, Deadline: deadline})
			saveTasksToFile("tasks.json")
			http.Redirect(w, r, "/tasks", http.StatusSeeOther) // Перенаправление на страницу задач
			return
		} else if index := r.FormValue("index"); index != "" {
			// Обработка изменения статуса задачи
			if idx, err := strconv.Atoi(index); err == nil && idx >= 0 && idx < len(tasks) {
				tasks[idx].Completed = !tasks[idx].Completed
				saveTasksToFile("tasks.json")
				http.Redirect(w, r, "/tasks", http.StatusSeeOther) // Перенаправление на страницу задач
				return
			}
		} else if deleteIndex := r.FormValue("deleteIndex"); deleteIndex != "" {
			// Обработка удаления задачи
			if idx, err := strconv.Atoi(deleteIndex); err == nil && idx >= 0 && idx < len(tasks) {
				tasks = append(tasks[:idx], tasks[idx+1:]...) // Удаляем задачу
				saveTasksToFile("tasks.json")
				http.Redirect(w, r, "/tasks", http.StatusSeeOther) // Перенаправление на страницу задач
				return
			}
		}
	}

	// Отображение задач
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>Туду Лист</title>
	<style>
		body {
			text-align: center;
			font-family: Arial, sans-serif;
			background-color: #f4f4f4; /* Цвет фона */
		}
		h1 {
			font-size: 2.5em;
			margin-bottom: 20px;
		}
		ul {
			list-style-type: none;
			padding: 0;
			margin: 0 auto; /* Центрирование списка */
			width: 50%; /* Ширина списка */
		}
		li {
			margin: 10px 0;
			font-size: 1.2em; /* Уменьшение размера шрифта */
			display: flex; /* Используем Flexbox для выравнивания */
			justify-content: space-between; /* Разделяем элементы по краям */
			align-items: center; /* Выравниваем по центру по вертикали */
		}
		label {
			flex: 1; /* Занимает оставшееся пространство */
			min-width: 200px; /* Минимальная ширина для метки */
			margin: 0; /* Убираем отступы для метки */
		}
		input[type="text"], input[type="date"] {
			font-size: 1.2em;
			padding: 10px;
			width: 250px;
		}
		button {
			font-size: 1em; /* Уменьшение размера шрифта для кнопок */
			padding: 8px 16px; /* Уменьшение отступов кнопок */
			cursor: pointer;
			margin-left: 5px; /* Уменьшение отступа между кнопками и метками */
		}
		.completed {
			color: green; /* Зеленый цвет для выполненной задачи */
			text-decoration: line-through; /* Зачеркивание выполненной задачи */
		}
		.completed-button {
			background-color: green; /* Зеленый цвет для кнопки выполненной задачи */
			color: white; /* Белый текст на кнопке */
		}
	</style>
</head>
<body>
	<h1>Список задач</h1>
	<div style="display: flex; justify-content: flex-start; margin-bottom: 30px; margin-left: 150px;"> <!-- Отступ влево -->
		<a href="/logout" style="font-size: 1.5em; padding: 10px 20px; background-color: #007BFF; color: white; text-decoration: none; border-radius: 5px;">Выйти</a>
	</div> <!-- Кнопка выхода -->
	<ul>
	{{range $index, $task := .}}
		<li>
			<label class="{{if $task.Completed}}completed{{end}}">{{$task.Title}} (Дедлайн: {{formatDate $task.Deadline}})</label>
			<form method="post" action="/tasks" style="display:inline;">
				<input type="hidden" name="index" value="{{$index}}">
				<button type="submit" class="{{if $task.Completed}}completed-button{{end}}">{{if $task.Completed}}Выполнена{{else}}Невыполнена{{end}}</button>
			</form>
			<form method="post" action="/tasks" style="display:inline;">
				<input type="hidden" name="deleteIndex" value="{{$index}}">
				<button type="submit">Удалить</button>
			</form>
		</li>
	{{end}}
	</ul>
	<form method="post" action="/tasks">
		<input type="text" name="title" placeholder="Введите название задачи" required>
		<input type="date" name="deadline" required>
		<button type="submit">Добавить задачу</button>
	</form>
</body>
</html>`

	// Создаем FuncMap для передачи функции форматирования даты в шаблон
	funcMap := template.FuncMap{
		"formatDate": formatDate,
	}

	t, err := template.New("tasks").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		fmt.Println("Ошибка при парсинге шаблона:", err)
		return
	}
	t.Execute(w, tasks)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Главная страница</title>
	</head>
	<body>
		<h1>Добро пожаловать!</h1>
		<p><a href="/login">Войти</a></p>
		<p><a href="/register">Зарегистрироваться</a></p>
	</body>
	</html>`
	w.Write([]byte(tmpl))
}

func StartServer() {
	const filename = "tasks.json"
	loadTasksFromFile(filename)

	http.HandleFunc("/glav", mainPageHandler)          // Главная страница
	http.HandleFunc("/register", auth.RegisterHandler) // Маршрут для регистрации
	http.HandleFunc("/login", auth.LoginHandler)       // Маршрут для входа
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/logout", auth.LogoutHandler) // Маршрут для выхода
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}

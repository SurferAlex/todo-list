package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"testi/auth"
	// Добавляем пакет time для работы с датами
)

type Task struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
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

}

// Обработчик для отображения новых задач и добавления новых задач
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		if title := r.FormValue("title"); title != "" {
			// Обработка добавления новой задачи
			tasks = append(tasks, Task{Title: title, Completed: false}) // Обновляем добавление задачи
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

	// Загружаем шаблон
	tmpl := template.New("tasks.html")
	tmpl, err := tmpl.ParseFiles("frontend/templates/tasks.html")
	if err != nil {
		fmt.Println("Ошибка при парсинге шаблона:", err)
		http.Error(w, "Ошибка при парсинге шаблона", http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон
	err = tmpl.Execute(w, tasks)
	if err != nil {
		fmt.Println("Ошибка при отображении задач:", err)
		http.Error(w, "Ошибка при отображении задач", http.StatusInternalServerError)
		return
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("frontend/templates/home.html"))
	tmpl.Execute(w, nil)
}

func StartServer() {
	const filename = "tasks.json"
	loadTasksFromFile(filename)

	http.HandleFunc("/home", mainPageHandler)          // Главная страница
	http.HandleFunc("/register", auth.RegisterHandler) // Маршрут для регистрации
	http.HandleFunc("/login", auth.LoginHandler)       // Маршрут для входа
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/logout", auth.LogoutHandler) // Маршрут для выхода
	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}

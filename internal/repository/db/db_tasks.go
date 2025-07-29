package db

import (
	"testi/internal/entity"
)

func InsertTask(task entity.Task) error {
	query := `INSERT INTO tasks(username, title, is_completed) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, task.Username, task.Title, task.Complete)
	return err
}

// Получаем задачи от пользователя
func GetTasksByUser(username string) ([]entity.Task, error) {
	rows, err := db.Query(`SELECT id, username, title, is_completed FROM tasks WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task
		err := rows.Scan(&t.ID, &t.Username, &t.Title, &t.Complete)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Переключение статуса выполнения задачи по ID
func ToggleCompleteByID(username string, id int) error {
	query := `UPDATE tasks SET is_completed = NOT is_completed WHERE id = $1 AND username = $2`
	_, err := db.Exec(query, id, username)
	return err
}

// Удаление задачи по ID
func DeleteTaskByID(username string, id int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND username = $2`
	_, err := db.Exec(query, id, username)
	return err
}

func CreateTasksTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			title TEXT NOT NULL,
			is_completed BOOLEAN DEFAULT FALSE
		)
	`)
	return err
}

package tasks

import (
	"encoding/json"

	"os"
	"sync"

	"github.com/google/uuid"

	"testi/internal/entity"
)

const tasksFile = "tasks.json"

type Repository struct {
	mu sync.Mutex
}

func NewRepoTasks() *Repository {
	return &Repository{}
}

func (r *Repository) loadAllTasks() ([]entity.Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []entity.Task{}, nil
		}
		return nil, err
	}
	var tasks []entity.Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func (r *Repository) saveAllTasks(tasks []entity.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, data, 0644)
}

// Добавление задачи
func (r *Repository) CreateTask(user entity.User, text string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.loadAllTasks()
	if err != nil {
		return err
	}

	newTask := entity.Task{
		ID:       uuid.NewString(),
		UserName: user.Username,
		Text:     text,
		Complete: false,
	}

	tasks = append(tasks, newTask)
	return r.saveAllTasks(tasks)
}

// Получить задачи по пользователю
func (r *Repository) GetTasksByUser(username string) ([]entity.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.loadAllTasks()
	if err != nil {
		return nil, err
	}

	var userTasks []entity.Task
	for _, task := range tasks {
		if task.UserName == username {
			userTasks = append(userTasks, task)
		}
	}
	return userTasks, nil
}

// Удаление по ID
func (r *Repository) DeleteTaskByID(username, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.loadAllTasks()
	if err != nil {
		return err
	}

	newTasks := make([]entity.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.ID == id && task.UserName == username {
			continue // удаляем
		}
		newTasks = append(newTasks, task)
	}
	return r.saveAllTasks(newTasks)
}

// Переключение статуса по ID
func (r *Repository) ToggleCompleteByID(username, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.loadAllTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id && task.UserName == username {
			tasks[i].Complete = !task.Complete
			break
		}
	}

	return r.saveAllTasks(tasks)
}

package db

import (
	"database/sql"
	"log"
	"testi/internal/entity"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

var (
	// Переменная для хранения полдключения к базе данных
	db *sql.DB
)

// InitDB инициализирует подключение к базе данных
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем подключение к базе данных
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Успешно подключено к базе данных:", dataSourceName)
}

func GetDB() *sql.DB {
	return db
}

// GetUserByUsername получает пользователя по имени
func GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	query := `SELECT username, password FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err // Ошибка при выполнении запроса
	}
	return &user, nil // Возвращаем найденного пользователя
}

func InsertUser(user entity.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := db.Exec(query, user.Username, user.Password)
	return err
}

// DeleteUser удаляет пользователя по имени
func DeleteUser(username string) error {
	query := `DELETE FROM users WHERE username = $1`
	_, err := db.Exec(query, username)
	return err
}

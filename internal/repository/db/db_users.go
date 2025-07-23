package db

import (
	"database/sql"
	"testi/internal/entity"
)

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

// Получение userID по username
func GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	return userID, err
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

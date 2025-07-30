package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Session struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Created  int64  `json:"created"`
}

// Генерация уникального session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CreateSession(userID int, username string) (string, error) {
	sessionID := generateSessionID()
	session := Session{
		UserID:   userID,
		Username: username,
		Created:  time.Now().Unix(),
	}

	// Сериализация в JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации сессии: %v", err)
	}

	// Сохранение в Redis с TTL 7 дней
	ctx := context.Background()
	err = rdb.Set(ctx, "session:"+sessionID, sessionData, 7*24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения сессии: %v", err)
	}

	return sessionID, nil
}

func GetSession(sessionID string) (*Session, error) {
	ctx := context.Background()
	data, err := rdb.Get(ctx, "session:"+sessionID).Result()
	if err != nil {
		return nil, fmt.Errorf("сессия не найдена: %v", err)
	}

	var session Session
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации сессии: %v", err)
	}

	return &session, nil
}

func DeleteSession(sessionID string) error {
	ctx := context.Background()
	err := rdb.Del(ctx, "session:"+sessionID).Err()
	if err != nil {
		return fmt.Errorf("ошибка удаления сессии: %v", err)
	}

	return nil
}

// Проверка существования сессии
func SessionExists(sessionID string) bool {
	ctx := context.Background()
	exists, err := rdb.Exists(ctx, "session:"+sessionID).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

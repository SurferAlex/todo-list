package session

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	// Проверка подключения
	ctx := context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}

	log.Println("Успешно подключено к Redis")
}

func GetRedisClient() *redis.Client {
	return rdb
}

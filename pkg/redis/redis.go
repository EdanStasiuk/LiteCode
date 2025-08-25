package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test connection
	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}
}

package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	pass := os.Getenv("REDIS_PASSWORD")
	db := 0

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
}

package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func RedisInit() (*redis.Client, context.Context) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	return rdb, ctx
}

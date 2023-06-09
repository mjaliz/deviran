package initializers

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         context.Context
)

func ConnectRedis(config *AppConfig) {
	Ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		panic(err)
	}

	err := RedisClient.Set(Ctx, "test", "How to Refresh Access Tokens the Right Way in Golang", 0).Err()
	if err != nil {
		panic(err)
	}

	log.Println("✅ Redis client connected successfully...")
}

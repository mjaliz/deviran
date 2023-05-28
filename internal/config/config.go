package config

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	RedisCtx    context.Context
}

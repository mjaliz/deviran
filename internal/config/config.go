package config

import "gorm.io/gorm"

type AppConfig struct {
	DB *gorm.DB
}

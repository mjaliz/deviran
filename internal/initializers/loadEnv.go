package initializers

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	RedisUri      string `mapstructure:"REDIS_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
}

var Config *AppConfig

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return err
	}
	Config = &AppConfig{
		DBHost:                 appConfig.DBHost,
		DBUserName:             appConfig.DBUserName,
		DBUserPassword:         appConfig.DBUserPassword,
		DBName:                 appConfig.DBName,
		DBPort:                 appConfig.DBPort,
		ServerPort:             appConfig.ServerPort,
		RedisUri:               appConfig.RedisUri,
		RedisPassword:          appConfig.RedisPassword,
		AccessTokenPrivateKey:  appConfig.AccessTokenPrivateKey,
		AccessTokenPublicKey:   appConfig.AccessTokenPublicKey,
		RefreshTokenPrivateKey: appConfig.RefreshTokenPrivateKey,
		RefreshTokenPublicKey:  appConfig.RefreshTokenPublicKey,
		AccessTokenExpiresIn:   appConfig.AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  appConfig.RefreshTokenExpiresIn,
		AccessTokenMaxAge:      appConfig.AccessTokenMaxAge,
		RefreshTokenMaxAge:     appConfig.RefreshTokenMaxAge,
	}
	return nil
}

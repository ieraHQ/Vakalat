package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string `mapstructure:"APP_NAME"`
		Env     string `mapstructure:"APP_ENV"`
		Port    string `mapstructure:"APP_PORT"`
		Debug   bool   `mapstructure:"APP_DEBUG"`
	}
	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     string `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		SSLMode  string `mapstructure:"DB_SSL_MODE"`
	}
	Redis struct {
		Host     string `mapstructure:"REDIS_HOST"`
		Port     string `mapstructure:"REDIS_PORT"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	}
	JWT struct {
		Secret          string `mapstructure:"JWT_SECRET"`
		ExpirationHours int    `mapstructure:"JWT_EXPIRATION_HOURS"`
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Override with environment variables
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	if config.JWT.Secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &config, nil
}
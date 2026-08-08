package config

import (
	"log"

	"github.com/spf13/viper"
)

// Every nested struct below carries `mapstructure:",squash"` so its fields
// decode from the same flat key namespace as top-level config (DB_HOST,
// AI_BASE_URL, ...) instead of a nested "database.db_host" path that nothing
// ever produces. Without squash, viper.Unmarshal silently leaves every one of
// these fields at its zero value regardless of env vars, a config file, or
// SetDefault — the previous version of this file only "worked" for
// JWT.Secret because that field had a separate manual os.Getenv override.
type Config struct {
	App struct {
		Name  string `mapstructure:"APP_NAME"`
		Env   string `mapstructure:"APP_ENV"`
		Port  string `mapstructure:"APP_PORT"`
		Debug bool   `mapstructure:"APP_DEBUG"`
	} `mapstructure:",squash"`
	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     string `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		SSLMode  string `mapstructure:"DB_SSL_MODE"`
	} `mapstructure:",squash"`
	Redis struct {
		Host     string `mapstructure:"REDIS_HOST"`
		Port     string `mapstructure:"REDIS_PORT"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	} `mapstructure:",squash"`
	JWT struct {
		Secret          string `mapstructure:"JWT_SECRET"`
		ExpirationHours int    `mapstructure:"JWT_EXPIRATION_HOURS"`
	} `mapstructure:",squash"`
	AI struct {
		BaseURL        string `mapstructure:"AI_BASE_URL"`
		Model          string `mapstructure:"AI_MODEL"`
		EmbeddingModel string `mapstructure:"AI_EMBEDDING_MODEL"`
	} `mapstructure:",squash"`
}

// envKeys are every key the Config struct expects to be populated from the
// environment. AutomaticEnv() alone isn't enough for Unmarshal to pick these
// up — viper only resolves a key from the environment if it already knows
// about the key (via SetDefault, an explicit Set, or BindEnv).
var envKeys = []string{
	"APP_NAME", "APP_ENV", "APP_PORT", "APP_DEBUG",
	"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
	"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
	"JWT_SECRET", "JWT_EXPIRATION_HOURS",
	"AI_BASE_URL", "AI_MODEL", "AI_EMBEDDING_MODEL",
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	for _, key := range envKeys {
		if err := viper.BindEnv(key); err != nil {
			return nil, err
		}
	}

	// Set default values
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)
	// Local LLM endpoint, OpenAI-compatible (e.g. Ollama, LM Studio).
	viper.SetDefault("AI_BASE_URL", "http://localhost:11434/v1")
	viper.SetDefault("AI_MODEL", "llama3.1")
	// 384-dimension embedding model to match the schema's VECTOR(384) columns.
	viper.SetDefault("AI_EMBEDDING_MODEL", "all-minilm")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.JWT.Secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &config, nil
}

package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server `mapstructure:"server"`
		DB     DB     `mapstructure:"db"`
		Auth   Auth   `mapstructure:"auth"`
	}

	Server struct {
		Port int `mapstructure:"port"`
	}

	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
		PoolMax  int    `mapstructure:"pool_max"`
	}

	Auth struct {
		SecretKey          string        `mapstructure:"secret_key"`
		AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
		RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
	}
)

func New() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()

	// Config file setup
	v.SetConfigName("configs")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// Enable env var overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values (optional, but recommended)
	v.SetDefault("server.port", 8080)
	v.SetDefault("db.sslmode", "disable")
	v.SetDefault("auth.access_token_expiry", "15m")
	v.SetDefault("auth.refresh_token_expiry", "24h") // 1 day

	// Explicit BindEnv for keys that won't match automatically
	v.BindEnv("auth.secret_key", "SECRET_KEY")
	v.BindEnv("db.user", "DATABASE_USER")
	v.BindEnv("db.password", "DATABASE_PASSWORD")
	v.BindEnv("db.host", "DATABASE_HOST")
	v.BindEnv("db.port", "DATABASE_PORT")
	v.BindEnv("db.name", "DATABASE_NAME")
	v.BindEnv("db.sslmode", "DATABASE_SSLMODE")

	// Read YAML file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	// Optional: validate required fields
	if cfg.Auth.SecretKey == "" {
		return nil, fmt.Errorf("auth.secret_key must be set via config or env")
	}

	return &cfg, nil
}

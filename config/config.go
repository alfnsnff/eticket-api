package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Configuration struct {
		Server Server `mapstructure:"server"`
		DB     DB     `mapstructure:"db"`
		Auth   Auth   `mapstructure:"auth"`
		Tripay Tripay `mapstructure:"tripay"`
	}

	Server struct {
		Port int `mapstructure:"port"`
	}

	Auth struct {
		SecretKey          string        `mapstructure:"secret_key"`
		AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
		RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
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

	Tripay struct {
		ApiKey        string `mapstructure:"api_key"`
		PrivateApiKey string `mapstructure:"private_api_key"`
		MerhcantCode  string `mapstructure:"merchant_code"`
	}
)

func New() (*Configuration, error) {
	_ = godotenv.Load()

	v := viper.New()

	// Enable env var overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values (optional, but recommended)
	v.SetDefault("server.port", 8080)
	v.SetDefault("db.sslmode", "disable")
	v.SetDefault("auth.access_token_expiry", "15m")
	v.SetDefault("auth.refresh_token_expiry", "24h")

	// Explicit BindEnv for keys that won't match automatically
	v.BindEnv("server.port", "PORT")
	v.BindEnv("auth.secret_key", "SECRET_KEY")
	v.BindEnv("db.user", "DATABASE_USER")
	v.BindEnv("db.password", "DATABASE_PASSWORD")
	v.BindEnv("db.host", "DATABASE_HOST")
	v.BindEnv("db.port", "DATABASE_PORT")
	v.BindEnv("db.name", "DATABASE_NAME")
	v.BindEnv("db.sslmode", "DATABASE_SSLMODE")
	v.BindEnv("tripay.api_key", "TRIPAY_API_KEY")
	v.BindEnv("tripay.private_api_key", "TRIPAY_PRIVATE_API_KEY")
	v.BindEnv("tripay.merchant_code", "TRIPAY_MERCHANT_CODE")

	// Unmarshal into struct
	var cfg Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	// Optional: validate required fields
	if cfg.Auth.SecretKey == "" {
		return nil, fmt.Errorf("auth.secret_key must be set via config or env")
	}

	return &cfg, nil
}

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
		Server     Server     `mapstructure:"server"`
		DB         DB         `mapstructure:"db"`
		Auth       Auth       `mapstructure:"auth"`
		Tripay     Tripay     `mapstructure:"tripay"`
		SMTPMailer SMTPMailer `mapstructure:"smtpmailer"`
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

	SMTPMailer struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
	}
)

func New() (*Configuration, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found or failed to load it:", err)
	}

	v := viper.New()

	// Enable env var overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables manually
	bindEnvs := map[string]string{
		"server.port":               "PORT",
		"auth.secret_key":           "SECRET_KEY",
		"auth.access_token_expiry":  "ACCESS_TOKEN_EXPIRY",
		"auth.refresh_token_expiry": "REFRESH_TOKEN_EXPIRY",

		"tripay.api_key":         "TRIPAY_API_KEY",
		"tripay.private_api_key": "TRIPAY_PRIVATE_API_KEY",
		"tripay.merchant_code":   "TRIPAY_MERCHANT_CODE",

		"db.host":     "DATABASE_HOST",
		"db.port":     "DATABASE_PORT",
		"db.name":     "DATABASE_NAME",
		"db.user":     "DATABASE_USER",
		"db.password": "DATABASE_PASSWORD",
		"db.sslmode":  "DATABASE_SSLMODE",

		"smtpmailer.host":     "MAILER_HOST",
		"smtpmailer.port":     "MAILER_PORT",
		"smtpmailer.from":     "MAILER_FROM",
		"smtpmailer.username": "MAILER_USERNAME",
		"smtpmailer.password": "MAILER_PASSWORD",
	}

	for key, env := range bindEnvs {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("failed to bind env var %s: %w", env, err)
		}
	}

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

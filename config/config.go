package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server `mapstructure:"server"`
		DB     DB     `mapstructure:"db"`
		Token  Token  `mapstructure:"Token"`
		Tripay Tripay `mapstructure:"tripay"`
		SMTP   SMTP   `mapstructure:"smtp"`
	}

	Server struct {
		Port int `mapstructure:"port"`
	}

	Token struct {
		SecretKey string `mapstructure:"secret_key"`
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

	SMTP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found or failed to load it:", err)
	}
	v := viper.New()
	// Enable env var overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables manually
	bindEnvs := map[string]string{
		"server.port":      "PORT",
		"token.secret_key": "SECRET_KEY",

		"tripay.api_key":         "TRIPAY_API_KEY",
		"tripay.private_api_key": "TRIPAY_PRIVATE_API_KEY",
		"tripay.merchant_code":   "TRIPAY_MERCHANT_CODE",

		"db.host":     "DATABASE_HOST",
		"db.port":     "DATABASE_PORT",
		"db.name":     "DATABASE_NAME",
		"db.user":     "DATABASE_USER",
		"db.password": "DATABASE_PASSWORD",
		"db.sslmode":  "DATABASE_SSLMODE",

		"smtp.host":     "MAILER_HOST",
		"smtp.port":     "MAILER_PORT",
		"smtp.from":     "MAILER_FROM",
		"smtp.username": "MAILER_USERNAME",
		"smtp.password": "MAILER_PASSWORD",
	}

	for key, env := range bindEnvs {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("failed to bind env var %s: %w", env, err)
		}
	}

	// Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	return &cfg, nil
}

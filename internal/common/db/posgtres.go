package db

import (
	"eticket-api/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Host, cfg.DB.Port, cfg.DB.SSLMode,
	)

	// Configure GORM with custom settings
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := pg.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to configure database pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(100)          // sesuai kapasitas PostgreSQL
	sqlDB.SetMaxIdleConns(10)           // idle pool cukup besar
	sqlDB.SetConnMaxLifetime(time.Hour) // koneksi akan refresh tiap 1 jam

	return pg, nil
}

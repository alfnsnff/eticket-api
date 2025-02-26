package main

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain"
	"eticket-api/pkg/db/postgres"
	"eticket-api/pkg/utils/conf"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configPath := conf.GetConf(os.Getenv("ENV"))

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate your models (creating tables, etc.)
	if err := db.AutoMigrate(&domain.Route{}, &domain.Class{}, &domain.Schedule{}, &domain.Ship{}, &domain.Harbor{}, &domain.Booking{}, &domain.Ticket{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Set up Gin router and routes
	router := gin.Default()
	route.Setup(router, db)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

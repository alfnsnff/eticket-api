package main

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain/entities"
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
		// Log a non-fatal message. This is expected in environments
		// where a .env file is not used (e.g., production in Azure).
		log.Printf("INFO: Could not load .env file. Assuming environment variables are set externally: %v", err)
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
	if err := db.AutoMigrate(&entities.Route{}, &entities.Class{}, &entities.Schedule{}, &entities.Ship{}, &entities.Harbor{}, &entities.Booking{}, &entities.Ticket{}, &entities.ShipClass{}, &entities.Price{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	// Set up Gin router and routes
	router := gin.New()
	route.Setup(router, db)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

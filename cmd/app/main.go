package main

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain/entity"
	"eticket-api/pkg/db/postgres"
	"eticket-api/pkg/utils/conf"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Log a non-fatal message. This is expected in environments
		// where a .env file is not used (e.g., production in Azure).
		log.Printf("INFO: Could not load .env file: %v", err)
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
	if err := db.AutoMigrate(&entity.Route{}, &entity.Class{}, &entity.Schedule{}, &entity.Ship{}, &entity.Harbor{}, &entity.Booking{}, &entity.Ticket{}, &entity.Manifest{}, &entity.Fare{}, &entity.Allocation{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Define allowed origins (your frontend URL)
	origins := []string{"http://localhost:3000", "https://tiket-hebat.vercel.app/"}

	// Apply CORS middleware FIRST
	// This ensures CORS headers are processed for all routes defined afterwards
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Include OPTIONS for preflight requests
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// These routes will now be protected by the CORS middleware
	route.Setup(router, db)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

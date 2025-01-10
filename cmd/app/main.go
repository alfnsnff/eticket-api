package main

import (
    "fmt"
    "log"
    "os"
    "eticket-api/internal/handler"
    "eticket-api/internal/repository"
    "eticket-api/internal/domain"
    "eticket-api/internal/service"
    "eticket-api/internal/transport"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Build connection string from environment variables
    dsn := fmt.Sprintf(
        "user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSL_MODE"),
    )

    // Connect to PostgreSQL using GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

   // Automatically migrate your models (creating tables, etc.)
   if err := db.AutoMigrate(&domain.Ticket{}); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
    }


    // Set up the application layers
    ticketRepo := &repository.TicketRepositoryImpl{DB: db}
    ticketService := &service.TicketService{Repo: ticketRepo}
    ticketHandler := &handler.TicketHandler{Service: ticketService}

    // Set up Gin router and routes
    router := gin.Default()
    transport.SetupRoutes(router, ticketHandler)

    // Run the server
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

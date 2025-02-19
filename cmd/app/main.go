package main

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/db/postgres"
	"eticket-api/pkg/utils"
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

	configPath := utils.GetConfEnv(os.Getenv("ENV"))

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
	if err := db.AutoMigrate(&domain.Route{}, &domain.Class{}, &domain.Ticket{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Set up the application layers
	ticketRepository := &repository.TicketRepository{DB: db}
	ticketUsecase := &usecase.TicketUsecase{TicketRepository: ticketRepository}
	ticketController := &controller.TicketController{TicketUsecase: *ticketUsecase}

	classRepository := &repository.ClassRepository{DB: db}
	classUsecase := &usecase.ClassUsecase{ClassRepository: classRepository}
	classController := &controller.ClassController{ClassUsecase: classUsecase}

	routeRepository := &repository.RouteRepository{DB: db}
	routeUsecase := &usecase.RouteUsecase{RouteRepository: routeRepository}
	routeController := &controller.RouteController{RouteUsecase: routeUsecase}

	// Set up Gin router and routes
	router := gin.Default()
	route.SetupRoutes(router, ticketController, classController, routeController)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

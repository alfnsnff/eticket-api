package main

import (
	"context"
	"eticket-api/config"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain/entity"
	authentity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/job"
	"eticket-api/pkg/db/postgres"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate your models (creating tables, etc.)
	if err := db.AutoMigrate(&entity.Route{},
		&entity.Class{},
		&entity.Schedule{},
		&entity.Ship{},
		&entity.Harbor{},
		&entity.Booking{},
		&entity.ClaimSession{},
		&entity.Ticket{},
		&entity.Manifest{},
		&entity.Fare{},
		&entity.Allocation{},
		&authentity.Role{},
		&authentity.User{},
		&authentity.UserRole{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	cleanupJob := job.SetupJob(db)

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // Ensure cancel is called

		log.Println("Starting cleanup job goroutine...")
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		if err := cleanupJob.Run(ctx); err != nil {
			log.Printf("Initial cleanup job run failed: %v", err)
		}

		for {
			select {
			case <-ticker.C:
				log.Println("Triggering scheduled cleanup job run...")
				if err := cleanupJob.Run(ctx); err != nil {
					log.Printf("Scheduled cleanup job run failed: %v", err)
				}
			case <-ctx.Done():
				log.Println("Cleanup job goroutine shutting down.")
				return
			}
		}
	}()

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	origins := []string{"http://localhost:3000", "https://tiket-hebat.vercel.app"}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Include OPTIONS for preflight requests
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.Setup(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

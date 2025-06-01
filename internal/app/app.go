package app

import (
	"context"
	"eticket-api/config"
	"eticket-api/internal/delivery/http/route"
	"eticket-api/internal/domain/entity"
	authentity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/injector"
	"eticket-api/internal/job"
	"eticket-api/pkg/casbinx"
	"eticket-api/pkg/db/postgres"
	"fmt"

	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	ic, err := injector.InitializeContainer(cfg)
	if err != nil {
		log.Fatalf("Failed initiate injection container: %v", err)
	}

	// Initialize database connection
	db, err := postgres.New(ic.Cfg)
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
		&authentity.RefreshToken{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	cleanupJob := job.Setup(ic)

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

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			allowed := map[string]bool{
				"https://hoppscotch.io":                                             true,
				"http://localhost:3000":                                             true,
				"https://tiket-hebat.vercel.app":                                    true,
				"https://tiket-hebat-ardians-projects-01d38d65.vercel.app":          true,
				"https://tiket-hebat-git-main-ardians-projects-01d38d65.vercel.app": true,
			}
			return allowed[origin]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(func(c *gin.Context) {
		fmt.Println("Origin:", c.Request.Header.Get("Origin"))
		c.Next()
	})

	casbinService := casbinx.NewCasbinService(ic.Enforcer)
	casbinx.Policies(casbinService)

	route.Setup(router, ic)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

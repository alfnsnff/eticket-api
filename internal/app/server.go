//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"eticket-api/config"
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/module"
	"net/http"

	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func NewHTTPClient() *http.Client {
	return &http.Client{}
}

type Server struct {
	config *config.Configuration
	app    *gin.Engine
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewHTTPClient,
		NewDatabase,
		NewEnforcer,

		// Config
		config.New,

		jwt.New,
		tx.New,

		mailer.NewSMTPMailer,

		// Modules
		module.NewClientModule,
		module.NewRepositoryModule,
		module.NewUsecaseModule,
		module.NewControllerModule,
		module.NewRouteModule,
		module.NewMiddlewareModule,
		module.NewJobModule,

		// App Server
		NewServer,
	)))
}

func NewServer(config *config.Configuration) *Server {

	db, err := NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate your models (creating tables, etc.)
	if err := db.AutoMigrate(
		&entity.Route{},
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
		&entity.Role{},
		&entity.User{},
		&entity.RefreshToken{},
		&entity.PasswordReset{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := gin.Default()

	Bootstrap(app, NewHTTPClient(), config, db)
	app.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// TODO: Make dynamic via cfg
			allowed := map[string]bool{
				"http://localhost:3000":          true,
				"https://tiket-hebat.vercel.app": true,
				"https://www.tikethebat.live":    true,
				"https://tripay.co.id/":          true,
			}
			return allowed[origin]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	NewEnforcer()

	Setup(route, app)

	// Start background job
	go Job(job)

	return &Server{
		config: config,
		app:    app,
	}
}

func Job(job *module.JobModule) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting cleanup job goroutine...")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	if err := job.CleanupJob.Run(ctx); err != nil {
		log.Printf("Initial cleanup failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println("Running scheduled cleanup job...")
			if err := job.CleanupJob.Run(ctx); err != nil {
				log.Printf("Scheduled cleanup failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("Cleanup job shutting down.")
			return
		}
	}
}

func Setup(route *module.RouteModule, app *gin.Engine) {
	app.Group("/api/v1")
	route.AuthRoute.Set(app)
	route.AllocationRoute.Set(app)
	route.BookingRoute.Set(app)
	route.ClassRoute.Set(app)
	route.FareRoute.Set(app)
	route.HarborRoute.Set(app)
	route.ManifestRoute.Set(app)
	route.RoleRoute.Set(app)
	route.Routeouter.Set(app)
	route.ScheduleRoute.Set(app)
	// route.SessionRoute.Set(app)
	route.ShipRoute.Set(app)
	route.TicketRoute.Set(app)
	route.UserRoute.Set(app)
	route.PaymentRoute.Set(app)
}

func (server Server) App() *gin.Engine {
	return server.app
}

func (server Server) Config() *config.Configuration {
	return server.config
}

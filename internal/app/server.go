//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"eticket-api/config"
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/tx"
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

func NewServer(
	config *config.Configuration,
	route *module.RouterModule,
	job *module.JobModule,
) *Server {
	gin.SetMode(gin.DebugMode)

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// TODO: Make dynamic via cfg
			allowed := map[string]bool{
				"http://localhost:3000": true,
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

func Setup(route *module.RouterModule, app *gin.Engine) {
	group := app.Group("/api/v1")
	route.AllocationRouter.Set(app, group)
	route.AuthRouter.Set(app, group)
	route.BookingRouter.Set(app, group)
	route.ClassRouter.Set(app, group)
	route.FareRouter.Set(app, group)
	route.HarborRouter.Set(app, group)
	route.ManifestRouter.Set(app, group)
	route.RoleRouter.Set(app, group)
	route.RouteRouter.Set(app, group)
	route.ScheduleRouter.Set(app, group)
	route.SessionRouter.Set(app, group)
	route.ShipRouter.Set(app, group)
	route.TicketRouter.Set(app, group)
	route.UserRouter.Set(app, group)
	route.PaymentRouter.Set(app, group)
}

func (server Server) App() *gin.Engine {
	return server.app
}

func (server Server) Config() *config.Configuration {
	return server.config
}

//go:build wireinject
// +build wireinject

package app

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http"
	"eticket-api/internal/domain"
	"eticket-api/internal/job"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type Server struct {
	app *gin.Engine
}

// Wire injector
func New(cfg *config.Config) (*Server, error) {
	panic(wire.Build(
		CommonSet,
		RepositorySet,
		ClientSet,
		UsecaseSet,
		JobSet,
		RouterSet,
		NewServer,
	))
}

// NewServer menerima semua dependency yang dibutuhkan, Wire akan mengisi otomatis
func NewServer(
	db *gorm.DB,
	router *http.Router,
	claimSessionJob *job.ClaimSessionJob,
	emailJob *job.EmailJob,
) (*Server, error) {
	gin.SetMode(gin.DebugMode)
	app := gin.Default()

	// Migrasi database
	if err := db.AutoMigrate(
		&domain.Role{},
		&domain.User{},
		&domain.Ship{},
		&domain.Harbor{},
		&domain.Class{},
		&domain.Schedule{},
		&domain.Booking{},
		&domain.Quota{},
		&domain.ClaimSession{},
		&domain.ClaimItem{},
		&domain.Ticket{},
		&domain.RefreshToken{},
		&domain.PasswordReset{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			allowed := map[string]bool{
				"http://localhost:3000":          true,
				"https://tiket-hebat.vercel.app": true,
				"https://www.tikethebat.live":    true,
				"https://tripay.co.id/":          true,
			}
			return allowed[origin]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := app.Group("/api")

	// Register langsung per grup
	router.RegisterMetrics(api)
	router.RegisterV1(api.Group("/v1"))
	router.RegisterV2(api.Group("/v2"))
	go claimSessionJob.CleanExpiredClaimSession()

	return &Server{app: app}, nil
}

func (server Server) App() *gin.Engine {
	return server.app
}

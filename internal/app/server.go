package app

import (
	"eticket-api/config"
	"eticket-api/internal/common/db"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/domain"
	"fmt"

	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	app *gin.Engine
}

func NewApp(cfg *config.Config) (*Server, error) {
	fmt.Println(">>> NewServer CALLED")
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	db, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Automatically migrate your models (creating tables, etc.)
	if err := db.AutoMigrate(
		&domain.Route{},
		&domain.Class{},
		&domain.Schedule{},
		&domain.Ship{},
		&domain.Harbor{},
		&domain.Booking{},
		&domain.ClaimSession{},
		&domain.Ticket{},
		&domain.Manifest{},
		&domain.Fare{},
		&domain.Allocation{},
		&domain.Role{},
		&domain.User{},
		&domain.RefreshToken{},
		&domain.PasswordReset{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	bootstrap := &Bootstrap{
		Config:   cfg,
		App:      app,
		DB:       db,
		Client:   httpclient.NewHTTPClient(cfg),
		Enforcer: enforcer.NewCasbinEnforcer(cfg),
		Token:    token.NewJWT(cfg),
		Mailer:   mailer.NewSMTP(cfg),
		Log:      logger.NewLogrusLogger("eticket-api"),
		Validate: validator.NewValidator(cfg),
	}

	if err := NewBootstrap(bootstrap); err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		app: app}, nil

}

func (server Server) App() *gin.Engine {
	return server.app
}

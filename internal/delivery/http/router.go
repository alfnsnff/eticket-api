package http

import (
	"eticket-api/internal/delivery/http/controller"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AdvancedRouter uses reflection for ultimate clean architecture
type Router struct {
	Engine       *gin.Engine
	Dependencies *Dependencies
}

// ControllerRegistry defines controller configuration
type ControllerRegistry struct {
	Name        string
	Constructor interface{} // Controller constructor function
	Usecase     interface{} // Associated usecase
}

// NewAdvancedRouter creates router with reflection-based registration
func NewRouter(engine *gin.Engine, dependencies *Dependencies) *Router {
	return &Router{
		Engine:       engine,
		Dependencies: dependencies,
	}
}

// Setup configures routes with auto-discovery pattern
func (router *Router) Setup() {
	router.SetupGlobalMiddleware()
	router.SetupSystemRoutes()
	router.SetupControllers()
}

// setupGlobalMiddleware registers global middleware
func (router *Router) SetupGlobalMiddleware() {
	router.Engine.Use(router.Dependencies.Middlewares.Logger.Set())
	router.Engine.Use(router.Dependencies.Middlewares.Recovery.Set())
}

// setupSystemRoutes registers system routes
func (router *Router) SetupSystemRoutes() {
	router.Engine.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	router.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// setupControllersWithReflection automatically registers controllers
func (router *Router) SetupControllers() {
	// Controller registry - auto-discovery pattern
	registry := []ControllerRegistry{
		{
			Name:        "Allocation",
			Constructor: controller.NewAllocationController,
			Usecase:     router.Dependencies.Usecases.Allocation,
		},
		{
			Name:        "Auth",
			Constructor: controller.NewAuthController,
			Usecase:     router.Dependencies.Usecases.Auth,
		},
		{
			Name:        "Booking",
			Constructor: controller.NewBookingController,
			Usecase:     router.Dependencies.Usecases.Booking,
		},
		{
			Name:        "ClaimSession",
			Constructor: controller.NewClaimSessionController,
			Usecase:     router.Dependencies.Usecases.ClaimSession,
		},
		{
			Name:        "Class",
			Constructor: controller.NewClassController,
			Usecase:     router.Dependencies.Usecases.Class,
		},
		{
			Name:        "Fare",
			Constructor: controller.NewFareController,
			Usecase:     router.Dependencies.Usecases.Fare,
		},
		{
			Name:        "Harbor",
			Constructor: controller.NewHarborController,
			Usecase:     router.Dependencies.Usecases.Harbor,
		},
		{
			Name:        "Manifest",
			Constructor: controller.NewManifestController,
			Usecase:     router.Dependencies.Usecases.Manifest,
		},
		{
			Name:        "Payment",
			Constructor: controller.NewPaymentController,
			Usecase:     router.Dependencies.Usecases.Payment,
		},
		{
			Name:        "Role",
			Constructor: controller.NewRoleController,
			Usecase:     router.Dependencies.Usecases.Role,
		},
		{
			Name:        "Route",
			Constructor: controller.NewRouteController,
			Usecase:     router.Dependencies.Usecases.Route,
		},
		{
			Name:        "Schedule",
			Constructor: controller.NewScheduleController,
			Usecase:     router.Dependencies.Usecases.Schedule,
		},
		{
			Name:        "Ship",
			Constructor: controller.NewShipController,
			Usecase:     router.Dependencies.Usecases.Ship,
		},
		{
			Name:        "Ticket",
			Constructor: controller.NewTicketController,
			Usecase:     router.Dependencies.Usecases.Ticket,
		},
		{
			Name:        "User",
			Constructor: controller.NewUserController,
			Usecase:     router.Dependencies.Usecases.User,
		},
	}

	// Auto-register controllers using reflection
	for _, ctrl := range registry {
		router.RegisterController(ctrl)
	}
}

// registerController uses reflection to call controller constructors
func (router *Router) RegisterController(ctrl ControllerRegistry) {
	constructorValue := reflect.ValueOf(ctrl.Constructor)

	// Prepare constructor arguments
	args := []reflect.Value{
		reflect.ValueOf(router.Engine),                                // router
		reflect.ValueOf(router.Dependencies.Logger),                   // logger
		reflect.ValueOf(router.Dependencies.Validator),                // validator
		reflect.ValueOf(ctrl.Usecase),                                 // usecase
		reflect.ValueOf(router.Dependencies.Middlewares.Authenticate), // authenticate
		reflect.ValueOf(router.Dependencies.Middlewares.Authorize),    // authorize
	}

	// Call constructor with reflection
	constructorValue.Call(args)
}

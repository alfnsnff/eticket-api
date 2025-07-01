package http

import (
	"eticket-api/internal/delivery/http/middleware"
	v1 "eticket-api/internal/delivery/http/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register untuk /metrics dan /health
func (r *Router) RegisterMetrics(group *gin.RouterGroup) {
	group.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	group.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Register untuk /v1
func (r *Router) RegisterV1(group *gin.RouterGroup) {
	group.Use(middleware.Logger(r.Logger))
	group.Use(middleware.Recovery(r.Logger))
	protected := group.Group("")
	protected.Use(middleware.Authenticate(r.TokenUtil))

	v1.NewQuotaController(group, protected, r.Logger, r.Validator, r.Quota)
	v1.NewAuthController(group, protected, r.Logger, r.Validator, r.Auth)
	v1.NewBookingController(group, protected, r.Logger, r.Validator, r.Booking)
	v1.NewClassController(group, protected, r.Logger, r.Validator, r.Class)
	v1.NewClaimSessionController(group, protected, r.Logger, r.Validator, r.ClaimSession)
	v1.NewHarborController(group, protected, r.Logger, r.Validator, r.Harbor)
	v1.NewPaymentController(group, protected, r.Logger, r.Validator, r.Payment)
	v1.NewRoleController(group, protected, r.Logger, r.Validator, r.Role)
	v1.NewScheduleController(group, protected, r.Logger, r.Validator, r.Schedule)
	v1.NewShipController(group, protected, r.Logger, r.Validator, r.Ship)
	v1.NewTicketController(group, protected, r.Logger, r.Validator, r.Ticket)
	v1.NewUserController(group, protected, r.Logger, r.Validator, r.User)
}

// Register untuk /v2 (future)
func (r *Router) RegisterV2(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "API v2 belum tersedia",
		})
	})
}

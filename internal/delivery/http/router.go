package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/usecase"

	// "eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	v1 "eticket-api/internal/delivery/http/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(
	router *gin.RouterGroup,
	token_util token.TokenUtil,
	log logger.Logger,
	validate validator.Validator,
	quota *usecase.QuotaUsecase,
	auth *usecase.AuthUsecase,
	booking *usecase.BookingUsecase,
	class *usecase.ClassUsecase,
	harbor *usecase.HarborUsecase,
	role *usecase.RoleUsecase,
	schedule *usecase.ScheduleUsecase,
	ship *usecase.ShipUsecase,
	ticket *usecase.TicketUsecase,
	user *usecase.UserUsecase,
	payment *usecase.PaymentUsecase,
	claim_session *usecase.ClaimSessionUsecase,
) {
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))

	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authMiddleware := middleware.Authenticate(token_util)

	// API v1
	v1Group := router.Group("/v1")
	v1Protected := Protected(v1Group, authMiddleware)

	v1.NewQuotaController(v1Group, v1Protected, log, validate, quota)
	v1.NewAuthController(v1Group, v1Protected, log, validate, auth)
	v1.NewBookingController(v1Group, v1Protected, log, validate, booking)
	v1.NewClassController(v1Group, v1Protected, log, validate, class)
	v1.NewClaimSessionController(v1Group, v1Protected, log, validate, claim_session)
	v1.NewHarborController(v1Group, v1Protected, log, validate, harbor)
	v1.NewPaymentController(v1Group, v1Protected, log, validate, payment)
	v1.NewRoleController(v1Group, v1Protected, log, validate, role)
	v1.NewScheduleController(v1Group, v1Protected, log, validate, schedule)
	v1.NewShipController(v1Group, v1Protected, log, validate, ship)
	v1.NewTicketController(v1Group, v1Protected, log, validate, ticket)
	v1.NewUserController(v1Group, v1Protected, log, validate, user)

}

func Protected(rg *gin.RouterGroup, middleware ...gin.HandlerFunc) *gin.RouterGroup {
	group := rg.Group("")
	group.Use(middleware...)
	return group
}

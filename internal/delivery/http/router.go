package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/usecase"

	// "eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(
	router *gin.RouterGroup,
	token_util token.TokenUtil,
	log logger.Logger,
	validate validator.Validator,
	// allocation *usecase.AllocationUsecase,
	quota *usecase.QuotaUsecase,
	auth *usecase.AuthUsecase,
	booking *usecase.BookingUsecase,
	class *usecase.ClassUsecase,
	// fare *usecase.FareUsecase,
	harbor *usecase.HarborUsecase,
	// manifest *usecase.ManifestUsecase,
	role *usecase.RoleUsecase,
	// route *usecase.RouteUsecase,
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

	V1 := router.Group("/v1")
	protected := V1.Group("")
	protected.Use(middleware.Authenticate(token_util))
	controller.NewQuotaController(V1, protected, log, validate, quota)
	controller.NewAuthController(V1, protected, log, validate, auth)
	controller.NewBookingController(V1, protected, log, validate, booking)
	controller.NewClassController(V1, protected, log, validate, class)
	controller.NewClaimSessionController(V1, protected, log, validate, claim_session)
	controller.NewHarborController(V1, protected, log, validate, harbor)
	controller.NewPaymentController(V1, protected, log, validate, payment)
	controller.NewRoleController(V1, protected, log, validate, role)
	controller.NewScheduleController(V1, protected, log, validate, schedule)
	controller.NewShipController(V1, protected, log, validate, ship)
	controller.NewTicketController(V1, protected, log, validate, ticket)
	controller.NewUserController(V1, protected, log, validate, user)
}

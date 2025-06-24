package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/usecase/allocation"
	"eticket-api/internal/usecase/auth"
	"eticket-api/internal/usecase/booking"
	"eticket-api/internal/usecase/claim_session"
	"eticket-api/internal/usecase/class"
	"eticket-api/internal/usecase/fare"
	"eticket-api/internal/usecase/harbor"
	"eticket-api/internal/usecase/manifest"
	"eticket-api/internal/usecase/payment"
	"eticket-api/internal/usecase/role"
	"eticket-api/internal/usecase/route"
	"eticket-api/internal/usecase/schedule"
	"eticket-api/internal/usecase/ship"
	"eticket-api/internal/usecase/ticket"
	"eticket-api/internal/usecase/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(
	router *gin.Engine,
	log logger.Logger,
	validate validator.Validator,
	allocation_usecase *allocation.AllocationUsecase,
	auth_usecase *auth.AuthUsecase,
	booking_usecase *booking.BookingUsecase,
	role_usecase *role.RoleUsecase,
	claimSession_usecase *claim_session.ClaimSessionUsecase,
	class_usecase *class.ClassUsecase,
	fare_usecase *fare.FareUsecase,
	harbor_usecase *harbor.HarborUsecase,
	manifest_usecase *manifest.ManifestUsecase,
	payment_usecase *payment.PaymentUsecase,
	route_usecase *route.RouteUsecase,
	schedule_usecase *schedule.ScheduleUsecase,
	ship_usecase *ship.ShipUsecase,
	ticket_usecase *ticket.TicketUsecase,
	user_usecase *user.UserUsecase,
	authenticate *middleware.AuthenticateMiddleware,
	authorize *middleware.AuthorizeMiddleware,
) {

	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))

	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	controller.NewAllocationController(router, log, validate, allocation_usecase, authenticate, authorize)
	controller.NewAuthController(router, log, validate, auth_usecase, authenticate, authorize)
	controller.NewBookingController(router, log, validate, booking_usecase, authenticate, authorize)
	controller.NewClaimSessionController(router, log, validate, claimSession_usecase, authenticate, authorize)
	controller.NewClassController(router, log, validate, class_usecase, authenticate, authorize)
	controller.NewFareController(router, log, validate, fare_usecase, authenticate, authorize)
	controller.NewHarborController(router, log, validate, harbor_usecase, authenticate, authorize)
	controller.NewManifestController(router, log, validate, manifest_usecase, authenticate, authorize)
	controller.NewPaymentController(router, log, validate, payment_usecase, authenticate, authorize)
	controller.NewRoleController(router, log, validate, role_usecase, authenticate, authorize)
	controller.NewRouteController(router, log, validate, route_usecase, authenticate, authorize)
	controller.NewScheduleController(router, log, validate, schedule_usecase, authenticate, authorize)
	controller.NewShipController(router, log, validate, ship_usecase, authenticate, authorize)
	controller.NewTicketController(router, log, validate, ticket_usecase, authenticate, authorize)
	controller.NewUserController(router, log, validate, user_usecase, authenticate, authorize)

}

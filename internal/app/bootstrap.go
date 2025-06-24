package app

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/job"
	"eticket-api/internal/repository"
	"eticket-api/internal/runner"
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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Bootstrap struct {
	Config   *config.Config
	App      *gin.Engine
	DB       *gorm.DB
	Client   *httpclient.HTTP
	Enforcer enforcer.Enforcer
	Validate validator.Validator
	Token    token.TokenUtil
	Mailer   mailer.Mailer
	Log      logger.Logger
}

func NewBootstrap(cf *Bootstrap) error {
	TripayClient := client.NewTripayClient(cf.Client, &cf.Config.Tripay)

	UserRepository := repository.NewUserRepository()
	RoleRepository := repository.NewRoleRepository()
	AuthRepository := repository.NewAuthRepository()
	ShipRepository := repository.NewShipRepository()
	RouteRepository := repository.NewRouteRepository()
	HarborRepository := repository.NewHarborRepository()
	ClassRepository := repository.NewClassRepository()
	FareRepository := repository.NewFareRepository()
	ManifestRepository := repository.NewManifestRepository()
	ScheduleRepository := repository.NewScheduleRepository()
	AllocationRepository := repository.NewAllocationRepository()
	TicketRepository := repository.NewTicketRepository()
	BookingRepository := repository.NewBookingRepository()
	SessionRepository := repository.NewSessionRepository()

	RoleUsecase := role.NewRoleUsecase(cf.DB, RoleRepository)
	UserUsecase := user.NewUserUsecase(cf.DB, UserRepository)
	AuthUsecase := auth.NewAuthUsecase(cf.DB, AuthRepository, UserRepository, cf.Mailer, cf.Token)
	ShipUsecase := ship.NewShipUsecase(cf.DB, ShipRepository)
	RouteUsecase := route.NewRouteUsecase(cf.DB, RouteRepository)
	HarborUsecase := harbor.NewHarborUsecase(cf.DB, HarborRepository)
	ClassUsecase := class.NewClassUsecase(cf.DB, ClassRepository)
	FareUsecase := fare.NewFareUsecase(cf.DB, FareRepository)
	ManifestUsecase := manifest.NewManifestUsecase(cf.DB, ManifestRepository)
	AllocationUsecase := allocation.NewAllocationUsecase(cf.DB, AllocationRepository, FareRepository)
	ScheduleUsecase := schedule.NewScheduleUsecase(cf.DB, AllocationRepository, ClassRepository, FareRepository, ManifestRepository, ShipRepository, ScheduleRepository, TicketRepository)
	TicketUsecase := ticket.NewTicketUsecase(cf.DB, TicketRepository, ScheduleRepository, ManifestRepository, FareRepository)
	BookingUsecase := booking.NewBookingUsecase(cf.DB, BookingRepository)
	ClaimSessionUsecase := claim_session.NewClaimSessionUsecase(cf.DB, SessionRepository, TicketRepository, ScheduleRepository, AllocationRepository, ManifestRepository, FareRepository, BookingRepository)
	PaymentUsecase := payment.NewPaymentUsecase(cf.DB, TripayClient, SessionRepository, BookingRepository, TicketRepository, cf.Mailer)

	// Initialize the enforcer
	AllcationController := controller.NewAllocationController(cf.Log, cf.Validate, AllocationUsecase)
	AuthController := controller.NewAuthController(cf.Log, cf.Validate, AuthUsecase)
	BookingController := controller.NewBookingController(cf.Log, cf.Validate, BookingUsecase)
	ClaimSessionController := controller.NewClaimSessionController(cf.Log, cf.Validate, ClaimSessionUsecase)
	ClassController := controller.NewClassController(cf.Log, cf.Validate, ClassUsecase)
	FareController := controller.NewFareController(cf.Log, cf.Validate, FareUsecase)
	HarborController := controller.NewHarborController(cf.Log, cf.Validate, HarborUsecase)
	ManifestController := controller.NewManifestController(cf.Log, cf.Validate, ManifestUsecase)
	PaymentController := controller.NewPaymentController(cf.Log, cf.Validate, PaymentUsecase)
	ScheduleController := controller.NewScheduleController(cf.Log, cf.Validate, ScheduleUsecase)
	RoleController := controller.NewRoleController(cf.Log, cf.Validate, RoleUsecase)
	RouteController := controller.NewRouteController(cf.Log, cf.Validate, RouteUsecase)
	ShipController := controller.NewShipController(cf.Log, cf.Validate, ShipUsecase)
	TicketController := controller.NewTicketController(cf.Log, cf.Validate, TicketUsecase)
	UserController := controller.NewUserController(cf.Log, cf.Validate, UserUsecase)

	AuthenticateMiddleware := middleware.NewAuthenticateMiddleware(cf.Token)
	AuthorizeMiddleware := middleware.NewAuthorizeMiddleware(cf.Enforcer)
	LoggerMiddleware := middleware.NewLoggerMiddleware(cf.Log)
	// RecoveryMiddleware := middleware.RecoveryMiddleware(cf.Log)

	router := http.NewRouter(cf.App, &http.Controllers{
		Allocation:   AllcationController,
		Auth:         AuthController,
		Booking:      BookingController,
		ClaimSession: ClaimSessionController,
		Class:        ClassController,
		Fare:         FareController,
		Harbor:       HarborController,
		Manifest:     ManifestController,
		Payment:      PaymentController,
		Schedule:     ScheduleController,
		Role:         RoleController,
		Route:        RouteController,
		Ship:         ShipController,
		Ticket:       TicketController,
		User:         UserController,
	}, &http.Middlewares{
		Authenticate: AuthenticateMiddleware,
		Authorize:    AuthorizeMiddleware,
		Log:          LoggerMiddleware,
	})

	router.Setup()

	cleanupJob := job.NewCleanupJob(cf.DB, TicketRepository, SessionRepository)
	cleanupRunner := runner.NewCleanupRunner(cleanupJob)
	cleanupRunner.Start()

	// // Create and setup router
	// router := http.NewRouter(cf.App, dependencies)
	// router.Setup() // This was missing!

	return nil
}

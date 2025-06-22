package app

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
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
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Bootstrap struct {
	Config   *config.Config
	App      *gin.Engine
	DB       *gorm.DB
	Client   *http.Client
	Enforcer enforcer.Enforcer
	Validate validator.Validator
	Token    token.TokenUtil
	Mailer   mailer.Mailer
	Log      logger.Logger
}

func NewBootstrap(cf *Bootstrap) error {
	// === Middleware ===
	Authenticate := middleware.NewAuthenticateMiddleware(cf.Token)
	Authorize := middleware.NewAuthorizeMiddleware(cf.Enforcer)

	// === Role Domain ===
	RoleRepository := repository.NewRoleRepository()
	RoleUsecase := role.NewRoleUsecase(cf.DB, RoleRepository)
	controller.NewRoleController(cf.App, cf.Log, cf.Validate, RoleUsecase, Authenticate, Authorize)

	// === User Domain ===
	UserRepository := repository.NewUserRepository()
	UserUsecase := user.NewUserUsecase(cf.DB, UserRepository)
	controller.NewUserController(cf.App, cf.Log, cf.Validate, UserUsecase, Authenticate, Authorize)

	// === Auth Domain ===
	AuthRepository := repository.NewAuthRepository()
	AuthUsecase := auth.NewAuthUsecase(cf.DB, AuthRepository, UserRepository, cf.Mailer, cf.Token)
	controller.NewAuthController(cf.App, cf.Log, cf.Validate, AuthUsecase, Authenticate, Authorize)

	// === Ship Domain ===
	ShipRepository := repository.NewShipRepository()
	ShipUsecase := ship.NewShipUsecase(cf.DB, ShipRepository)
	controller.NewShipController(cf.App, cf.Log, cf.Validate, ShipUsecase, Authenticate, Authorize)

	// === Route Domain ===
	RouteRepository := repository.NewRouteRepository()
	RouteUsecase := route.NewRouteUsecase(cf.DB, RouteRepository)
	controller.NewRouteController(cf.App, cf.Log, cf.Validate, RouteUsecase, Authenticate, Authorize)

	// === Harbor Domain ===
	HarborRepository := repository.NewHarborRepository()
	HarborUsecase := harbor.NewHarborUsecase(cf.DB, HarborRepository)
	controller.NewHarborController(cf.App, cf.Log, cf.Validate, HarborUsecase, Authenticate, Authorize)

	// === Class Domain ===
	ClassRepository := repository.NewClassRepository(cf.DB)
	ClassUsecase := class.NewClassUsecase(cf.DB, ClassRepository)
	controller.NewClassController(cf.App, cf.Log, cf.Validate, ClassUsecase, Authenticate, Authorize)

	// === Fare Domain ===
	FareRepository := repository.NewFareRepository()
	FareUsecase := fare.NewFareUsecase(cf.DB, FareRepository)
	controller.NewFareController(cf.App, cf.Log, cf.Validate, FareUsecase, Authenticate, Authorize)

	// === Manifest Domain ===
	ManifestRepository := repository.NewManifestRepository()
	ManifestUsecase := manifest.NewManifestUsecase(cf.DB, ManifestRepository)
	controller.NewManifestController(cf.App, cf.Log, cf.Validate, ManifestUsecase, Authenticate, Authorize)

	// === Allocation Domain ===
	AllocationRepository := repository.NewAllocationRepository()
	AllocationUsecase := allocation.NewAllocationUsecase(cf.DB, AllocationRepository, FareRepository)
	controller.NewAllocationController(cf.App, cf.Log, cf.Validate, AllocationUsecase, Authenticate, Authorize)

	// === Ticket Domain ===
	TicketRepository := repository.NewTicketRepository()
	TicketUsecase := ticket.NewTicketUsecase(cf.DB, TicketRepository)
	controller.NewTicketController(cf.App, cf.Log, cf.Validate, TicketUsecase, Authenticate, Authorize)

	// === Schedule Domain ===
	ScheduleRepository := repository.NewScheduleRepository()
	ScheduleUsecase := schedule.NewScheduleUsecase(
		cf.DB,
		AllocationRepository,
		ClassRepository,
		FareRepository,
		ManifestRepository,
		ShipRepository,
		ScheduleRepository,
		TicketRepository,
	)
	controller.NewScheduleController(cf.App, cf.Log, cf.Validate, ScheduleUsecase, Authenticate, Authorize)

	// === Booking Domain ===
	BookingRepository := repository.NewBookingRepository()
	BookingUsecase := booking.NewBookingUsecase(cf.DB, BookingRepository)
	controller.NewBookingController(cf.App, cf.Log, cf.Validate, BookingUsecase, Authenticate, Authorize)

	// === Session Domain ===
	// === ClaimSession Domain ===
	SessionRepository := repository.NewSessionRepository()
	ClaimSessionUsecase := claim_session.NewClaimSessionUsecase(
		cf.DB,
		SessionRepository,
		TicketRepository,
		ScheduleRepository,
		AllocationRepository,
		ManifestRepository,
		FareRepository,
		BookingRepository,
	)
	controller.NewClaimSessionController(cf.App, cf.Log, cf.Validate, ClaimSessionUsecase, Authenticate, Authorize)
	// === Payment Domain ===
	TripayClient := client.NewTripayClient(cf.Client, &cf.Config.Tripay)
	PaymentUsecase := payment.NewPaymentUsecase(cf.DB, TripayClient, SessionRepository, BookingRepository, TicketRepository, cf.Mailer)
	controller.NewPaymentController(cf.App, cf.Log, cf.Validate, PaymentUsecase, Authenticate, Authorize)

	// === Cleanup Job ===
	cleanupJob := job.NewCleanupJob(cf.DB, TicketRepository, SessionRepository)
	cleanupRunner := runner.NewCleanupRunner(cleanupJob)
	cleanupRunner.Start()

	return nil
}

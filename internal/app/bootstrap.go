package app

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/repository"
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

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Bootstrap(gin *gin.Engine, http *http.Client, config *config.Configuration, db *gorm.DB, tm *jwt.TokenUtil, e *casbin.Enforcer, mailer *mailer.SMTPMailer) {
	// === Middleware ===
	Authenticate := middleware.NewAuthenticateMiddleware(tm)
	Authorize := middleware.NewAuthorizeMiddleware(e)

	// === Role Domain ===
	RoleRepository := repository.NewRoleRepository()
	RoleUsecase := role.NewRoleUsecase(db, RoleRepository)
	controller.NewRoleController(gin, RoleUsecase, Authenticate, Authorize)

	// === User Domain ===
	UserRepository := repository.NewUserRepository()
	UserUsecase := user.NewUserUsecase(db, UserRepository)
	controller.NewUserController(gin, UserUsecase, Authenticate, Authorize)

	// === Auth Domain ===
	AuthRepository := repository.NewAuthRepository()
	AuthUsecase := auth.NewAuthUsecase(db, AuthRepository, UserRepository, mailer, tm)
	controller.NewAuthController(gin, AuthUsecase, Authenticate, Authorize)

	// === Ship Domain ===
	ShipRepository := repository.NewShipRepository()
	ShipUsecase := ship.NewShipUsecase(db, ShipRepository)
	controller.NewShipController(gin, ShipUsecase, Authenticate, Authorize)

	// === Route Domain ===
	RouteRepository := repository.NewRouteRepository()
	RouteUsecase := route.NewRouteUsecase(db, RouteRepository)
	controller.NewRouteController(gin, RouteUsecase, Authenticate, Authorize)

	// === Harbor Domain ===
	HarborRepository := repository.NewHarborRepository()
	HarborUsecase := harbor.NewHarborUsecase(db, HarborRepository)
	controller.NewHarborController(gin, HarborUsecase, Authenticate, Authorize)

	// === Class Domain ===
	ClassRepository := repository.NewClassRepository(db)
	ClassUsecase := class.NewClassUsecase(db, ClassRepository)
	controller.NewClassController(gin, ClassUsecase, Authenticate, Authorize)

	// === Fare Domain ===
	FareRepository := repository.NewFareRepository()
	FareUsecase := fare.NewFareUsecase(db, FareRepository)
	controller.NewFareController(gin, FareUsecase, Authenticate, Authorize)

	// === Manifest Domain ===
	ManifestRepository := repository.NewManifestRepository()
	ManifestUsecase := manifest.NewManifestUsecase(db, ManifestRepository)
	controller.NewManifestController(gin, ManifestUsecase, Authenticate, Authorize)

	// === Allocation Domain ===
	AllocationRepository := repository.NewAllocationRepository()
	AllocationUsecase := allocation.NewAllocationUsecase(db, AllocationRepository, FareRepository)
	controller.NewAllocationController(gin, AllocationUsecase, Authenticate, Authorize)

	// === Ticket Domain ===
	TicketRepository := repository.NewTicketRepository()
	TicketUsecase := ticket.NewTicketUsecase(db, TicketRepository)
	controller.NewTicketController(gin, TicketUsecase, Authenticate, Authorize)

	// === Schedule Domain ===
	ScheduleRepository := repository.NewScheduleRepository()
	ScheduleUsecase := schedule.NewScheduleUsecase(
		db,
		AllocationRepository,
		ClassRepository,
		FareRepository,
		ManifestRepository,
		ShipRepository,
		ScheduleRepository,
		TicketRepository,
	)
	controller.NewScheduleController(gin, ScheduleUsecase, Authenticate, Authorize)

	// === Booking Domain ===
	BookingRepository := repository.NewBookingRepository()
	BookingUsecase := booking.NewBookingUsecase(db, BookingRepository)
	controller.NewBookingController(gin, BookingUsecase, Authenticate, Authorize)

	// === Session Domain ===
	// === ClaimSession Domain ===
	SessionRepository := repository.NewSessionRepository()
	ClaimSessionUsecase := claim_session.NewClaimSessionUsecase(
		db,
		SessionRepository,
		TicketRepository,
		ScheduleRepository,
		AllocationRepository,
		ManifestRepository,
		FareRepository,
		BookingRepository,
	)
	controller.NewClaimSessionController(gin, ClaimSessionUsecase, Authenticate, Authorize)

	// === Payment Domain ===
	TripayClient := client.NewTripayClient(http, &config.Tripay)
	PaymentUsecase := payment.NewPaymentUsecase(db, TripayClient, BookingRepository, TicketRepository, mailer)
	controller.NewPaymentController(gin, PaymentUsecase, Authenticate, Authorize)
}

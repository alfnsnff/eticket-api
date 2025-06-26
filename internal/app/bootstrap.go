package app

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/transactor"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http"
	"eticket-api/internal/job"
	"eticket-api/internal/repository"
	"eticket-api/internal/runner"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Bootstrap struct {
	Config     *config.Config
	App        *gin.Engine
	DB         *gorm.DB
	Transactor *transactor.Transactor
	Client     *httpclient.HTTP
	Enforcer   enforcer.Enforcer
	Validate   validator.Validator
	Token      token.TokenUtil
	Mailer     mailer.Mailer
	Log        logger.Logger
}

func NewBootstrap(bootstrap *Bootstrap) error {
	// === Role Domain ===
	RoleRepository := repository.NewRoleRepository()
	RoleUsecase := usecase.NewRoleUsecase(
		bootstrap.DB,
		RoleRepository,
	)

	// === User Domain ===
	UserRepository := repository.NewUserRepository()
	UserUsecase := usecase.NewUserUsecase(
		bootstrap.DB,
		UserRepository,
	)

	// === Auth Domain ===
	AuthRepository := repository.NewAuthRepository()
	AuthUsecase := usecase.NewAuthUsecase(
		bootstrap.DB,
		AuthRepository,
		UserRepository,
		bootstrap.Mailer,
		bootstrap.Token,
	)

	// === Ship Domain ===
	ShipRepository := repository.NewShipRepository()
	ShipUsecase := usecase.NewShipUsecase(
		bootstrap.DB,
		ShipRepository,
	)

	// // === Route Domain ===
	// RouteRepository := repository.NewRouteRepository()
	// RouteUsecase := usecase.NewRouteUsecase(
	// 	bootstrap.DB,
	// 	RouteRepository,
	// )

	// === Harbor Domain ===
	HarborRepository := repository.NewHarborRepository()
	HarborUsecase := usecase.NewHarborUsecase(
		bootstrap.DB,
		HarborRepository,
	)

	// === Class Domain ===
	ClassRepository := repository.NewClassRepository()
	ClassUsecase := usecase.NewClassUsecase(
		bootstrap.DB,
		ClassRepository,
	)

	// // === Fare Domain ===
	// FareRepository := repository.NewFareRepository()
	// FareUsecase := usecase.NewFareUsecase(
	// 	bootstrap.DB,
	// 	FareRepository,
	// )

	// // === Manifest Domain ===
	// ManifestRepository := repository.NewManifestRepository()
	// ManifestUsecase := usecase.NewManifestUsecase(
	// 	bootstrap.DB,
	// 	ManifestRepository,
	// )

	// // === Allocation Domain ===
	// AllocationRepository := repository.NewAllocationRepository()
	// AllocationUsecase := usecase.NewAllocationUsecase(
	// 	bootstrap.DB,
	// 	AllocationRepository,
	// 	FareRepository,
	// )

	QuotaRepository := repository.NewQuotaRepository()
	QuotaUsecase := usecase.NewQuotaUsecase(
		bootstrap.DB,
		QuotaRepository,
	)

	// === Schedule Domain ===
	// === Ticket Domain ===
	TicketRepository := repository.NewTicketRepository()
	ScheduleRepository := repository.NewScheduleRepository()
	ScheduleUsecase := usecase.NewScheduleUsecase(
		bootstrap.DB,
		// AllocationRepository,
		ClassRepository,
		// FareRepository,
		// ManifestRepository,
		ShipRepository,
		ScheduleRepository,
		TicketRepository,
	)
	TicketUsecase := usecase.NewTicketUsecase(
		bootstrap.DB,
		TicketRepository,
		ScheduleRepository,
		// ManifestRepository,
		// FareRepository,
	)

	// === Booking Domain ===
	BookingRepository := repository.NewBookingRepository()
	BookingUsecase := usecase.NewBookingUsecase(
		bootstrap.DB,
		BookingRepository,
	)

	// === ClaimSession Domain ===
	ClaimItemRepository := repository.NewClaimItemRepository()
	_ = usecase.NewClaimItemUsecase(
		bootstrap.DB,
		ClaimItemRepository,
	)

	// === ClaimSession Domain ===
	ClaimSessionRepository := repository.NewClaimSessionRepository()
	ClaimSessionUsecase := usecase.NewClaimSessionUsecase(
		bootstrap.DB,
		ClaimSessionRepository,
		ClaimItemRepository,
		TicketRepository,
		ScheduleRepository,
		// AllocationRepository,
		// ManifestRepository,
		// FareRepository,
		BookingRepository,
		QuotaRepository,
	)

	// === Payment Domain ===
	TripayClient := client.NewTripayClient(
		bootstrap.Client,
		&bootstrap.Config.Tripay,
	)
	PaymentUsecase := usecase.NewPaymentUsecase(
		bootstrap.DB,
		TripayClient,
		ClaimSessionRepository,
		BookingRepository,
		TicketRepository,
		bootstrap.Mailer,
	)

	// === Initialize HTTP Router ===
	api := bootstrap.App.Group("/api")
	http.NewRouter(
		api,
		bootstrap.Token,
		bootstrap.Log,
		bootstrap.Validate,
		// AllocationUsecase,
		QuotaUsecase,
		AuthUsecase,
		BookingUsecase,
		ClassUsecase,
		// FareUsecase,
		HarborUsecase,
		// ManifestUsecase,
		RoleUsecase,
		// RouteUsecase,
		ScheduleUsecase,
		ShipUsecase,
		TicketUsecase,
		UserUsecase,
		PaymentUsecase,
		ClaimSessionUsecase,
	)

	// === Initialize Cleanup Job ===
	cleanupJob := job.NewCleanupJob(
		bootstrap.DB,
		TicketRepository,
		ClaimSessionRepository,
	)
	cleanupRunner := runner.NewCleanupRunner(cleanupJob)
	cleanupRunner.Start()

	return nil
}

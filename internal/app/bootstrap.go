package app

import (
	"eticket-api/config"
	"eticket-api/internal/client"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http"
	"eticket-api/internal/job"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Bootstrap struct {
	Config     *config.Config
	App        *gin.Engine
	DB         *gorm.DB
	Transactor *transact.Transactor
	Client     *httpclient.HTTP
	Enforcer   enforcer.Enforcer
	Validate   validator.Validator
	Token      token.TokenUtil
	Mailer     mailer.Mailer
	Log        logger.Logger
}

func NewBootstrap(i *Bootstrap) error {
	// === Role Domain ===
	RoleRepository := repository.NewRoleRepository(i.DB)
	RoleUsecase := usecase.NewRoleUsecase(
		i.Transactor,
		RoleRepository,
	)

	// === User Domain ===
	UserRepository := repository.NewUserRepository(i.DB)
	UserUsecase := usecase.NewUserUsecase(
		i.Transactor,
		UserRepository,
	)

	// === Auth Domain ===
	AuthRepository := repository.NewAuthRepository(i.DB)
	AuthUsecase := usecase.NewAuthUsecase(
		i.Transactor,
		AuthRepository,
		UserRepository,
		i.Mailer,
		i.Token,
	)

	// === Ship Domain ===
	ShipRepository := repository.NewShipRepository(i.DB)
	ShipUsecase := usecase.NewShipUsecase(
		i.Transactor,
		ShipRepository,
	)

	// === Harbor Domain ===
	HarborRepository := repository.NewHarborRepository(i.DB)
	HarborUsecase := usecase.NewHarborUsecase(
		i.Transactor,
		HarborRepository,
	)

	// === Class Domain ===
	ClassRepository := repository.NewClassRepository(i.DB)
	ClassUsecase := usecase.NewClassUsecase(
		i.Transactor,
		ClassRepository,
	)

	QuotaRepository := repository.NewQuotaRepository(i.DB)
	QuotaUsecase := usecase.NewQuotaUsecase(
		i.Transactor,
		QuotaRepository,
	)

	// === Schedule Domain ===
	// === Ticket Domain ===
	TicketRepository := repository.NewTicketRepository(i.DB)
	ScheduleRepository := repository.NewScheduleRepository(i.DB)
	ScheduleUsecase := usecase.NewScheduleUsecase(
		i.Transactor,
		ClassRepository,
		ShipRepository,
		ScheduleRepository,
		TicketRepository,
	)
	TicketUsecase := usecase.NewTicketUsecase(
		i.Transactor,
		TicketRepository,
		ScheduleRepository,
		QuotaRepository,
	)

	// === Booking Domain ===
	BookingRepository := repository.NewBookingRepository(i.DB)
	BookingUsecase := usecase.NewBookingUsecase(
		i.Transactor,
		BookingRepository,
	)

	// === ClaimSession Domain ===
	ClaimItemRepository := repository.NewClaimItemRepository(i.DB)
	_ = usecase.NewClaimItemUsecase(
		i.Transactor,
		ClaimItemRepository,
	)

	// === Payment Domain ===
	TripayClient := client.NewTripayClient(
		i.Client,
		&i.Config.Tripay,
	)

	// === ClaimSession Domain ===
	ClaimSessionRepository := repository.NewClaimSessionRepository(i.DB)
	ClaimSessionUsecase := usecase.NewClaimSessionUsecase(
		i.Transactor,
		ClaimSessionRepository,
		ClaimItemRepository,
		TicketRepository,
		ScheduleRepository,
		BookingRepository,
		QuotaRepository,
		TripayClient,
	)
	ClaimSessionJob := job.NewClaimSessionJob(i.Log, ClaimSessionUsecase)
	ClaimSessionJob.CleanExpiredClaimSession()

	PaymentUsecase := usecase.NewPaymentUsecase(
		i.Transactor,
		TripayClient,
		BookingRepository,
		TicketRepository,
		QuotaRepository,
		i.Mailer,
	)

	// === Initialize HTTP Router ===
	api := i.App.Group("/api")
	http.NewRouter(
		api,
		i.Token,
		i.Log,
		i.Validate,

		QuotaUsecase,
		AuthUsecase,
		BookingUsecase,
		ClassUsecase,

		HarborUsecase,

		RoleUsecase,

		ScheduleUsecase,
		ShipUsecase,
		TicketUsecase,
		UserUsecase,
		PaymentUsecase,
		ClaimSessionUsecase,
	)

	// // === Initialize Cleanup Job ===
	// cleanupJob := job.NewCleanupJob(
	// 	i.Transactor,
	// 	ClaimSessionRepository,
	// )
	// cleanupRunner := runner.NewCleanupRunner(cleanupJob)
	// cleanupRunner.Start()

	return nil
}

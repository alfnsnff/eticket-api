package app

import (
	"eticket-api/internal/client"
	"eticket-api/internal/common/db"
	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/common/httpclient"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http"
	"eticket-api/internal/domain"
	"eticket-api/internal/job"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/google/wire"
)

var CommonSet = wire.NewSet(
	db.NewPostgres, // returns *gorm.DB

	transact.NewTransactionManager,
	wire.Bind(new(transact.Transactor), new(*transact.Gotann)), // ⬅️ ini WAJIB

	enforcer.NewCasbinEnforcer, // returns *casbin.Enforcer

	logger.NewLogrus, // returns Logger              // returns Mailer            // returns TokenUtil
	wire.Value("eticket-api"),

	validator.NewValidator, // returns Validator

	token.NewJWT,
	wire.Bind(new(token.TokenUtil), new(*token.JWT)), // ✅ add this

	mailer.NewSMTP,
	wire.Bind(new(mailer.Mailer), new(*mailer.SMTP)), // ✅ add this

	// ✅ HTTP Client
	httpclient.NewHTTPClient, // <-- You need this to get *httpclient.HTTP

)

var RepositorySet = wire.NewSet(
	// --- Concrete constructors ---
	repository.NewRoleRepository,
	repository.NewUserRepository,
	repository.NewRefreshTokenRepository,
	repository.NewShipRepository,
	repository.NewHarborRepository,
	repository.NewClassRepository,
	repository.NewQuotaRepository,
	repository.NewScheduleRepository,
	repository.NewTicketRepository,
	repository.NewBookingRepository,
	repository.NewClaimItemRepository,
	repository.NewClaimSessionRepository,

	// --- Interface bindings ---
	wire.Bind(new(domain.RoleRepository), new(*repository.RoleRepository)),
	wire.Bind(new(domain.UserRepository), new(*repository.UserRepository)),
	wire.Bind(new(domain.RefreshTokenRepository), new(*repository.RefreshTokenRepository)),
	wire.Bind(new(domain.ShipRepository), new(*repository.ShipRepository)),
	wire.Bind(new(domain.HarborRepository), new(*repository.HarborRepository)),
	wire.Bind(new(domain.ClassRepository), new(*repository.ClassRepository)),
	wire.Bind(new(domain.QuotaRepository), new(*repository.QuotaRepository)),
	wire.Bind(new(domain.ScheduleRepository), new(*repository.ScheduleRepository)),
	wire.Bind(new(domain.TicketRepository), new(*repository.TicketRepository)),
	wire.Bind(new(domain.BookingRepository), new(*repository.BookingRepository)),
	wire.Bind(new(domain.ClaimItemRepository), new(*repository.ClaimItemRepository)),
	wire.Bind(new(domain.ClaimSessionRepository), new(*repository.ClaimSessionRepository)),
)

var ClientSet = wire.NewSet(
	client.NewTripayClient,

	wire.Bind(new(domain.TripayClient), new(*client.TripayClient)), // ✅ add this
	// ...dst
)

var UsecaseSet = wire.NewSet(
	usecase.NewRoleUsecase,
	usecase.NewUserUsecase,
	usecase.NewAuthUsecase,
	usecase.NewShipUsecase,
	usecase.NewHarborUsecase,
	usecase.NewClassUsecase,
	usecase.NewQuotaUsecase,
	usecase.NewScheduleUsecase,
	usecase.NewTicketUsecase,
	usecase.NewBookingUsecase,
	usecase.NewClaimItemUsecase,
	usecase.NewClaimSessionUsecase,
	usecase.NewPaymentUsecase,
	// ...dst
)

var JobSet = wire.NewSet(
	job.NewClaimSessionJob,
	// job.NewEmailJobQueue, // <--- tambahkan ini
)

var RouterSet = wire.NewSet(
	http.NewRouter,
)

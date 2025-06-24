package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
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
)

// Dependencies contains all router dependencies
type Dependencies struct {
	Logger      logger.Logger
	Validator   validator.Validator
	Middlewares *MiddlewareDependencies
	Usecases    *UsecaseDependencies
}

// MiddlewareDeps groups all middleware dependencies
type MiddlewareDependencies struct {
	Authenticate *middleware.AuthenticateMiddleware
	Logger       *middleware.LoggerMiddleware
	Recovery     *middleware.RecoveryMiddleware
	Authorize    *middleware.AuthorizeMiddleware
}

// UsecaseDeps groups all usecase dependencies
type UsecaseDependencies struct {
	Allocation   *allocation.AllocationUsecase
	Auth         *auth.AuthUsecase
	Booking      *booking.BookingUsecase
	Role         *role.RoleUsecase
	ClaimSession *claim_session.ClaimSessionUsecase
	Class        *class.ClassUsecase
	Fare         *fare.FareUsecase
	Harbor       *harbor.HarborUsecase
	Manifest     *manifest.ManifestUsecase
	Payment      *payment.PaymentUsecase
	Route        *route.RouteUsecase
	Schedule     *schedule.ScheduleUsecase
	Ship         *ship.ShipUsecase
	Ticket       *ticket.TicketUsecase
	User         *user.UserUsecase
}

// NewDependencies creates a new dependencies container
func NewDependencies(
	log logger.Logger,
	validate validator.Validator,
	middleware *MiddlewareDependencies,
	usecases *UsecaseDependencies,
) *Dependencies {
	return &Dependencies{
		Logger:      log,
		Validator:   validate,
		Middlewares: middleware,
		Usecases:    usecases,
	}
}

// NewMiddlewareDeps creates middleware dependencies container
func NewMiddlewareDependencies(
	authenticate *middleware.AuthenticateMiddleware,
	logger *middleware.LoggerMiddleware,
	recovery *middleware.RecoveryMiddleware,
	authorize *middleware.AuthorizeMiddleware,
) *MiddlewareDependencies {
	return &MiddlewareDependencies{
		Authenticate: authenticate,
		Logger:       logger,
		Recovery:     recovery,
		Authorize:    authorize,
	}
}

// NewUsecaseDeps creates usecase dependencies container
func NewUsecaseDependencies(
	allocation *allocation.AllocationUsecase,
	auth *auth.AuthUsecase,
	booking *booking.BookingUsecase,
	role *role.RoleUsecase,
	claimSession *claim_session.ClaimSessionUsecase,
	class *class.ClassUsecase,
	fare *fare.FareUsecase,
	harbor *harbor.HarborUsecase,
	manifest *manifest.ManifestUsecase,
	payment *payment.PaymentUsecase,
	route *route.RouteUsecase,
	schedule *schedule.ScheduleUsecase,
	ship *ship.ShipUsecase,
	ticket *ticket.TicketUsecase,
	user *user.UserUsecase,
) *UsecaseDependencies {
	return &UsecaseDependencies{
		Allocation:   allocation,
		Auth:         auth,
		Booking:      booking,
		Role:         role,
		ClaimSession: claimSession,
		Class:        class,
		Fare:         fare,
		Harbor:       harbor,
		Manifest:     manifest,
		Payment:      payment,
		Route:        route,
		Schedule:     schedule,
		Ship:         ship,
		Ticket:       ticket,
		User:         user,
	}
}

package http

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/usecase"
)

type Router struct {
	TokenUtil token.TokenUtil
	Logger    logger.Logger
	Validator validator.Validator

	Quota        *usecase.QuotaUsecase
	Auth         *usecase.AuthUsecase
	Booking      *usecase.BookingUsecase
	Class        *usecase.ClassUsecase
	Harbor       *usecase.HarborUsecase
	Role         *usecase.RoleUsecase
	Schedule     *usecase.ScheduleUsecase
	Ship         *usecase.ShipUsecase
	Ticket       *usecase.TicketUsecase
	User         *usecase.UserUsecase
	Payment      *usecase.PaymentUsecase
	ClaimSession *usecase.ClaimSessionUsecase
}

// NewRouter is Wire-compatible constructor
func NewRouter(
	tokenUtil token.TokenUtil,
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
	claimSession *usecase.ClaimSessionUsecase,
) *Router {
	return &Router{
		TokenUtil:    tokenUtil,
		Logger:       log,
		Validator:    validate,
		Quota:        quota,
		Auth:         auth,
		Booking:      booking,
		Class:        class,
		Harbor:       harbor,
		Role:         role,
		Schedule:     schedule,
		Ship:         ship,
		Ticket:       ticket,
		User:         user,
		Payment:      payment,
		ClaimSession: claimSession,
	}
}

package module

import (
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/tx"
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

type UsecaseModule struct {
	AuthUsecase *auth.AuthUsecase
	RoleUsecase *role.RoleUsecase
	UserUsecase *user.UserUsecase

	ShipUsecase       *ship.ShipUsecase
	AllocationUsecase *allocation.AllocationUsecase
	ManifestUsecase   *manifest.ManifestUsecase
	TicketUsecase     *ticket.TicketUsecase
	FareUsecase       *fare.FareUsecase
	ScheduleUsecase   *schedule.ScheduleUsecase
	BookingUsecase    *booking.BookingUsecase
	SessionUsecase    *claim_session.SessionUsecase
	RouteUsecase      *route.RouteUsecase
	HarborUsecase     *harbor.HarborUsecase
	ClassUsecase      *class.ClassUsecase
	PaymentUsecase    *payment.PaymentUsecase
}

func NewUsecaseModule(tx *tx.TxManager, repository *RepositoryModule, client *ClientModule, tm *jwt.TokenUtil, mailer *mailer.SMTPMailer) *UsecaseModule {
	return &UsecaseModule{

		AuthUsecase: auth.NewAuthUsecase(tx, repository.AuthRepository, repository.UserRepository, mailer, tm),
		RoleUsecase: role.NewRoleUsecase(tx, repository.RoleRepository),
		UserUsecase: user.NewUserUsecase(tx, repository.UserRepository),

		ShipUsecase:       ship.NewShipUsecase(tx, repository.ShipRepository),
		AllocationUsecase: allocation.NewAllocationUsecase(tx, repository.AllocationRepository, repository.ScheduleRepository, repository.FareRepository),
		ManifestUsecase:   manifest.NewManifestUsecase(tx, repository.ManifestRepository),
		TicketUsecase:     ticket.NewTicketUsecase(tx, repository.TicketRepository, repository.ScheduleRepository, repository.FareRepository, repository.SessionRepository),
		FareUsecase:       fare.NewFareUsecase(tx, repository.FareRepository),
		ScheduleUsecase:   schedule.NewScheduleUsecase(tx, repository.AllocationRepository, repository.ClassRepository, repository.FareRepository, repository.ManifestRepository, repository.RouteRepository, repository.ShipRepository, repository.ScheduleRepository, repository.TicketRepository),
		BookingUsecase:    booking.NewBookingUsecase(tx, repository.BookingRepository, repository.TicketRepository, repository.SessionRepository),
		SessionUsecase:    claim_session.NewSessionUsecase(tx, repository.SessionRepository, repository.TicketRepository, repository.ScheduleRepository, repository.AllocationRepository, repository.ManifestRepository, repository.FareRepository, repository.BookingRepository, client.TripayClient),
		RouteUsecase:      route.NewRouteUsecase(tx, repository.RouteRepository),
		HarborUsecase:     harbor.NewHarborUsecase(tx, repository.HarborRepository),
		ClassUsecase:      class.NewClassUsecase(tx, repository.ClassRepository),
		PaymentUsecase:    payment.NewPaymentUsecase(tx, client.TripayClient, repository.BookingRepository, repository.TicketRepository),
	}
}

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

	"gorm.io/gorm"
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

func NewUsecaseModule(tx *tx.TxManager, db *gorm.DB, repository *RepositoryModule, client *ClientModule, tm *jwt.TokenUtil, mailer *mailer.SMTPMailer) *UsecaseModule {
	return &UsecaseModule{

		AuthUsecase: auth.NewAuthUsecase(tx, db, repository.AuthRepository, repository.UserRepository, mailer, tm),
		RoleUsecase: role.NewRoleUsecase(db, repository.RoleRepository),
		UserUsecase: user.NewUserUsecase(db, repository.UserRepository),

		ShipUsecase:       ship.NewShipUsecase(db, repository.ShipRepository),
		AllocationUsecase: allocation.NewAllocationUsecase(tx, db, repository.AllocationRepository, repository.ScheduleRepository, repository.FareRepository),
		ManifestUsecase:   manifest.NewManifestUsecase(db, repository.ManifestRepository),
		TicketUsecase:     ticket.NewTicketUsecase(db, repository.TicketRepository, repository.ScheduleRepository, repository.FareRepository, repository.SessionRepository),
		FareUsecase:       fare.NewFareUsecase(db, repository.FareRepository),
		ScheduleUsecase:   schedule.NewScheduleUsecase(tx, db, repository.AllocationRepository, repository.ClassRepository, repository.FareRepository, repository.ManifestRepository, repository.RouteRepository, repository.ShipRepository, repository.ScheduleRepository, repository.TicketRepository),
		BookingUsecase:    booking.NewBookingUsecase(tx, db, repository.BookingRepository, repository.TicketRepository, repository.SessionRepository),
		SessionUsecase:    claim_session.NewSessionUsecase(tx, repository.SessionRepository, repository.TicketRepository, repository.ScheduleRepository, repository.AllocationRepository, repository.ManifestRepository, repository.FareRepository, repository.BookingRepository, client.TripayClient),
		RouteUsecase:      route.NewRouteUsecase(db, repository.RouteRepository),
		HarborUsecase:     harbor.NewHarborUsecase(db, repository.HarborRepository),
		ClassUsecase:      class.NewClassUsecase(db, repository.ClassRepository),
		PaymentUsecase:    payment.NewPaymentUsecase(tx, db, client.TripayClient, repository.BookingRepository, repository.TicketRepository, mailer),
	}
}

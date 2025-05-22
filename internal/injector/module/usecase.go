package module

import (
	"eticket-api/internal/usecase"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/tx"
)

type UsecaseModule struct {
	AuthUsecase     *authusecase.AuthUsecase
	RoleUsecase     *authusecase.RoleUsecase
	UserUsecase     *authusecase.UserUsecase
	UserRoleUsecase *authusecase.UserRoleUsecase

	ShipUsecase       *usecase.ShipUsecase
	AllocationUsecase *usecase.AllocationUsecase
	ManifestUsecase   *usecase.ManifestUsecase
	TicketUsecase     *usecase.TicketUsecase
	FareUsecase       *usecase.FareUsecase
	ScheduleUsecase   *usecase.ScheduleUsecase
	BookingUsecase    *usecase.BookingUsecase
	SessionUsecase    *usecase.SessionUsecase
	RouteUsecase      *usecase.RouteUsecase
	HarborUsecase     *usecase.HarborUsecase
	ClassUsecase      *usecase.ClassUsecase
}

func NewUsecaseModule(tx *tx.TxManager, repository *RepositoryModule, tm *jwt.TokenManager) *UsecaseModule {
	return &UsecaseModule{

		AuthUsecase:     authusecase.NewAuthUsecase(tx, repository.AuthRepository, repository.UserRepository, tm),
		RoleUsecase:     authusecase.NewRoleUsecase(tx, repository.RoleRepository),
		UserUsecase:     authusecase.NewUserUsecase(tx, repository.UserRepository),
		UserRoleUsecase: authusecase.NewUserRoleUsecase(tx, repository.RoleRepository, repository.UserRepository, repository.UserRoleRepository),

		ShipUsecase:       usecase.NewShipUsecase(tx, repository.ShipRepository),
		AllocationUsecase: usecase.NewAllocationUsecase(tx, repository.AllocationRepository, repository.ScheduleRepository, repository.FareRepository),
		ManifestUsecase:   usecase.NewManifestUsecase(tx, repository.ManifestRepository),
		TicketUsecase:     usecase.NewTicketUsecase(tx, repository.TicketRepository, repository.ScheduleRepository, repository.FareRepository, repository.SessionRepository),
		FareUsecase:       usecase.NewFareUsecase(tx, repository.FareRepository),
		ScheduleUsecase:   usecase.NewScheduleUsecase(tx, repository.AllocationRepository, repository.ClassRepository, repository.FareRepository, repository.ManifestRepository, repository.RouteRepository, repository.ShipRepository, repository.ScheduleRepository, repository.TicketRepository),
		BookingUsecase:    usecase.NewBookingUsecase(tx, repository.BookingRepository, repository.TicketRepository, repository.SessionRepository),
		SessionUsecase:    usecase.NewSessionUsecase(tx, repository.SessionRepository, repository.TicketRepository, repository.ScheduleRepository, repository.AllocationRepository, repository.ManifestRepository, repository.FareRepository),
		RouteUsecase:      usecase.NewRouteUsecase(tx, repository.RouteRepository),
		HarborUsecase:     usecase.NewHarborUsecase(tx, repository.HarborRepository),
		ClassUsecase:      usecase.NewClassUsecase(tx, repository.ClassRepository),
	}
}

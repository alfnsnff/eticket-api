package module

import (
	"eticket-api/internal/repository"
	authrepo "eticket-api/internal/repository/auth"
)

type RepositoryModule struct {
	AuthRepository     *authrepo.AuthRepository
	RoleRepository     *authrepo.RoleRepository
	UserRepository     *authrepo.UserRepository
	UserRoleRepository *authrepo.UserRoleRepository

	ShipRepository       *repository.ShipRepository
	AllocationRepository *repository.AllocationRepository
	ManifestRepository   *repository.ManifestRepository
	TicketRepository     *repository.TicketRepository
	FareRepository       *repository.FareRepository
	ScheduleRepository   *repository.ScheduleRepository
	BookingRepository    *repository.BookingRepository
	SessionRepository    *repository.SessionRepository
	RouteRepository      *repository.RouteRepository
	HarborRepository     *repository.HarborRepository
	ClassRepository      *repository.ClassRepository
}

func NewRepositoryModule() *RepositoryModule {
	return &RepositoryModule{
		AuthRepository:     authrepo.NewAuthRepository(),
		RoleRepository:     authrepo.NewRoleRepository(),
		UserRepository:     authrepo.NewUserRepository(),
		UserRoleRepository: authrepo.NewUserRoleRepository(),

		ShipRepository:       repository.NewShipRepository(),
		AllocationRepository: repository.NewAllocationRepository(),
		ManifestRepository:   repository.NewManifestRepository(),
		TicketRepository:     repository.NewTicketRepository(),
		FareRepository:       repository.NewFareRepository(),
		ScheduleRepository:   repository.NewScheduleRepository(),
		BookingRepository:    repository.NewBookingRepository(),
		SessionRepository:    repository.NewSessionRepository(),
		RouteRepository:      repository.NewRouteRepository(),
		HarborRepository:     repository.NewHarborRepository(),
		ClassRepository:      repository.NewClassRepository(),
	}
}

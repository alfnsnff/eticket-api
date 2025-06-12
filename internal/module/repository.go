package module

import (
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type RepositoryModule struct {
	AuthRepository *repository.AuthRepository
	RoleRepository *repository.RoleRepository
	UserRepository *repository.UserRepository

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

func NewRepositoryModule(db *gorm.DB) *RepositoryModule {
	return &RepositoryModule{
		AuthRepository: repository.NewAuthRepository(),
		RoleRepository: repository.NewRoleRepository(),
		UserRepository: repository.NewUserRepository(),

		TicketRepository:     repository.NewTicketRepository(),
		SessionRepository:    repository.NewSessionRepository(),
		ShipRepository:       repository.NewShipRepository(),
		AllocationRepository: repository.NewAllocationRepository(),
		ManifestRepository:   repository.NewManifestRepository(),
		FareRepository:       repository.NewFareRepository(),
		ScheduleRepository:   repository.NewScheduleRepository(),
		BookingRepository:    repository.NewBookingRepository(),
		RouteRepository:      repository.NewRouteRepository(),
		HarborRepository:     repository.NewHarborRepository(),
		ClassRepository:      repository.NewClassRepository(db),
	}
}

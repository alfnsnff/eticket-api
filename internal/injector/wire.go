//go:build wireinject
// +build wireinject

package injector

import (
	"eticket-api/config"
	"eticket-api/internal/delivery/http/controller"
	authcontroller "eticket-api/internal/delivery/http/controller/auth"

	"eticket-api/internal/repository"
	authrepository "eticket-api/internal/repository/auth"

	"eticket-api/internal/usecase"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/db/postgres"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/tx"

	"github.com/google/wire"
)

func InitializeContainer(cfg *config.Config) (*Container, error) {
	wire.Build(
		postgres.New,
		jwt.NewTokenManager,

		// TxManager binding
		tx.NewGormTxManager,
		wire.Bind(new(tx.TxManager), new(*tx.GormTxManager)),

		// Auth Repos
		authrepository.NewAuthRepository,
		authrepository.NewRoleRepository,
		authrepository.NewUserRepository,
		authrepository.NewUserRoleRepository,

		// App Repos
		repository.NewShipRepository,
		repository.NewAllocationRepository,
		repository.NewManifestRepository,
		repository.NewTicketRepository,
		repository.NewFareRepository,
		repository.NewScheduleRepository,
		repository.NewBookingRepository,
		repository.NewSessionRepository,
		repository.NewRouteRepository,
		repository.NewHarborRepository,
		repository.NewClassRepository,

		// Usecases
		authusecase.NewAuthUsecase,
		authusecase.NewRoleUsecase,
		authusecase.NewUserUsecase,
		authusecase.NewUserRoleUsecase,

		usecase.NewShipUsecase,
		usecase.NewAllocationUsecase,
		usecase.NewManifestUsecase,
		usecase.NewTicketUsecase,
		usecase.NewFareUsecase,
		usecase.NewScheduleUsecase,
		usecase.NewBookingUsecase,
		usecase.NewSessionUsecase,
		usecase.NewRouteUsecase,
		usecase.NewHarborUsecase,
		usecase.NewClassUsecase,

		// Controllers
		authcontroller.NewAuthController,
		authcontroller.NewRoleController,
		authcontroller.NewUserController,
		authcontroller.NewUserRoleController,

		controller.NewShipController,
		controller.NewAllocationController,
		controller.NewManifestController,
		controller.NewTicketController,
		controller.NewFareController,
		controller.NewScheduleController,
		controller.NewBookingController,
		controller.NewSessionController,
		controller.NewRouteController,
		controller.NewHarborController,
		controller.NewClassController,

		NewContainer,
	)
	return &Container{}, nil
}

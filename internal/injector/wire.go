//go:build wireinject
// +build wireinject

package injector

import (
	"eticket-api/config"
	"eticket-api/internal/injector/module"
	"eticket-api/pkg/casbinx"
	"eticket-api/pkg/db/postgres"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/tx"

	"github.com/google/wire"
)

func InitializeContainer(cfg *config.Config) (*Container, error) {
	wire.Build(
		// Core dependencies
		postgres.New,
		jwt.New,
		tx.New,
		casbinx.NewEnforcer,
		// Your internal module wiring
		module.NewRepositoryModule,
		module.NewUsecaseModule,
		module.NewControllerModule,

		// Final application container
		NewContainer,
	)
	return &Container{}, nil
}

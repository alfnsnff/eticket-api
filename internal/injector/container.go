package injector

import (
	"eticket-api/config"
	"eticket-api/internal/injector/module"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/tx"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type Container struct {
	Cfg          *config.Config
	DB           *gorm.DB
	Tx           *tx.TxManager
	TokenManager *jwt.TokenManager
	Enforcer     *casbin.Enforcer
	Repository   *module.RepositoryModule
	Usecase      *module.UsecaseModule
	Controller   *module.ControllerModule
}

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	tx *tx.TxManager,
	tm *jwt.TokenManager,
	enforcer *casbin.Enforcer,
	repository *module.RepositoryModule,
	usecase *module.UsecaseModule,
	controller *module.ControllerModule,
) *Container {
	return &Container{
		Cfg:          cfg,
		DB:           db,
		Tx:           tx,
		TokenManager: tm,
		Enforcer:     enforcer,
		Repository:   repository,
		Usecase:      usecase,
		Controller:   controller,
	}
}

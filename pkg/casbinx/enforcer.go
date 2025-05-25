package casbinx

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewEnforcer(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("failed to create adapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		log.Fatalf("failed to create enforcer: %v", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy: %v", err)
	}

	return enforcer
}

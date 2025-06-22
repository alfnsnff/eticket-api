package enforcer

import (
	"eticket-api/config"
	"log"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type Enforcer interface {
	Enforce(sub, obj, act string) (bool, error)
}

type CasbinEnforcer struct {
	enf *casbin.Enforcer
}

func NewCasbinEnforcer(cfg *config.Config) Enforcer {
	adapter := fileadapter.NewAdapter("config/policy.csv")

	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		log.Fatalf("failed to create enforcer: %v", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy: %v", err)
	}

	return &CasbinEnforcer{enf: enforcer}
}

func (c *CasbinEnforcer) Enforce(sub, obj, act string) (bool, error) {
	return c.enf.Enforce(sub, obj, act)
}

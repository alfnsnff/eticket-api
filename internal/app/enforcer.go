package app

import (
	"log"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func NewEnforcer() *casbin.Enforcer {
	// Create file adapter for policy.csv
	adapter := fileadapter.NewAdapter("config/policy.csv")

	// Load model from file and use adapter
	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		log.Fatalf("failed to create enforcer: %v", err)
	}

	// Load the policy from file into memory
	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy: %v", err)
	}

	return enforcer
}

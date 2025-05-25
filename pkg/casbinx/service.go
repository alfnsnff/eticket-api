package casbinx

import (
	"github.com/casbin/casbin/v2"
)

type CasbinService struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinService(enforcer *casbin.Enforcer) *CasbinService {
	return &CasbinService{Enforcer: enforcer}
}

func (s *CasbinService) AddPermission(role, obj, act string) (bool, error) {
	return s.Enforcer.AddPolicy(role, obj, act)
}

func (s *CasbinService) RemovePermission(role, obj, act string) (bool, error) {
	return s.Enforcer.RemovePolicy(role, obj, act)
}

func (s *CasbinService) GetPermissions(role string) ([][]string, error) {
	return s.Enforcer.GetFilteredPolicy(0, role)
}

func (s *CasbinService) ClearPermissions(role string) (bool, error) {
	return s.Enforcer.RemoveFilteredPolicy(0, role)
}

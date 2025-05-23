package casbin

import "github.com/casbin/casbin/v2"

type Service struct {
	e *casbin.Enforcer
}

func NewService(e *casbin.Enforcer) *Service {
	return &Service{e: e}
}

func (s *Service) AssignRole(userID, role string) error {
	_, err := s.e.AddGroupingPolicy(userID, role)
	return err
}

func (s *Service) AddPermission(role, path, method string) error {
	_, err := s.e.AddPolicy(role, path, method)
	return err
}

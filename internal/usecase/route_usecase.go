package usecase

import (
	"errors"
	"eticket-api/internal/domain"
	"fmt"
)

type RouteUsecase struct {
	RouteRepository domain.RouteRepository
}

func NewRouteUsecase(routeRepository domain.RouteRepository) RouteUsecase {
	return RouteUsecase{RouteRepository: routeRepository}
}

// CreateRoute validates and creates a new Route
func (s *RouteUsecase) CreateRoute(route *domain.Route) error {
	if route.DepartureHarbor == "" || route.ArrivalHarbor == "" {
		return fmt.Errorf("route name cannot be empty")
	}
	return s.RouteRepository.Create(route)
}

// GetAllRoutees retrieves all Routees
func (s *RouteUsecase) GetAllRoutes() ([]*domain.Route, error) {
	return s.RouteRepository.GetAll()
}

// GetRouteByID retrieves a Route by its ID
func (s *RouteUsecase) GetRouteByID(id uint) (*domain.Route, error) {
	Route, err := s.RouteRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if Route == nil {
		return nil, errors.New("route not found")
	}
	return Route, nil
}

// UpdateRoute updates an existing Route
func (s *RouteUsecase) UpdateRoute(route *domain.Route) error {
	if route.ID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if route.DepartureHarbor == "" || route.ArrivalHarbor == "" {
		return fmt.Errorf("route name cannot be empty")
	}
	return s.RouteRepository.Update(route)
}

// DeleteRoute deletes a Route by its ID
func (s *RouteUsecase) DeleteRoute(id uint) error {
	Route, err := s.RouteRepository.GetByID(id)
	if err != nil {
		return err
	}
	if Route == nil {
		return errors.New("route not found")
	}
	return s.RouteRepository.Delete(id)
}

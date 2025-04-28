package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	"fmt"
)

type RouteUsecase struct {
	RouteRepository *repository.RouteRepository
}

func NewRouteUsecase(routeRepository *repository.RouteRepository) RouteUsecase {
	return RouteUsecase{RouteRepository: routeRepository}
}

// CreateRoute validates and creates a new Route
func (s *RouteUsecase) CreateRoute(route *entities.Route) error {
	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}
	return s.RouteRepository.Create(route)
}

// GetAllRoutees retrieves all Routees
func (s *RouteUsecase) GetAllRoutes() ([]*entities.Route, error) {
	return s.RouteRepository.GetAll()
}

// GetRouteByID retrieves a Route by its ID
func (s *RouteUsecase) GetRouteByID(id uint) (*entities.Route, error) {
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
func (s *RouteUsecase) UpdateRoute(id uint, route *entities.Route) error {
	route.ID = id

	if route.ID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
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

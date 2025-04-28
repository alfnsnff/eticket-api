package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type RouteUsecase struct {
	DB              *gorm.DB
	RouteRepository *repository.RouteRepository
}

func NewRouteUsecase(db *gorm.DB, routeRepository *repository.RouteRepository) *RouteUsecase {
	return &RouteUsecase{
		DB:              db,
		RouteRepository: routeRepository,
	}
}

// GetAllRoutes retrieves all routes
func (s *RouteUsecase) GetAllRoutes(ctx context.Context) ([]*entities.Route, error) {
	var routes []*entities.Route

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		routes, err = s.RouteRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all routes: %w", err)
	}

	return routes, nil
}

// CreateRoute validates and creates a new Route
func (s *RouteUsecase) CreateRoute(ctx context.Context, route *entities.Route) error {
	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.RouteRepository.Create(txDB, route)
	})
}

// GetRouteByID retrieves a Route by its ID
func (s *RouteUsecase) GetRouteByID(ctx context.Context, id uint) (*entities.Route, error) {
	var route *entities.Route

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		route, err = s.RouteRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get route by ID: %w", err)
	}

	if route == nil {
		return nil, errors.New("route not found")
	}

	return route, nil
}

// UpdateRoute updates an existing Route
func (s *RouteUsecase) UpdateRoute(ctx context.Context, id uint, route *entities.Route) error {
	if id == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	route.ID = id

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.RouteRepository.Update(txDB, route)
	})
}

// DeleteRoute deletes a Route by its ID
func (s *RouteUsecase) DeleteRoute(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		route, err := s.RouteRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if route == nil {
			return errors.New("route not found")
		}
		return s.RouteRepository.Delete(txDB, id)
	})
}

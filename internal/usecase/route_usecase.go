package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type RouteUsecase struct {
	DB              *gorm.DB
	RouteRepository *repository.RouteRepository
}

func NewRouteUsecase(db *gorm.DB, route_repository *repository.RouteRepository) *RouteUsecase {
	return &RouteUsecase{
		DB:              db,
		RouteRepository: route_repository,
	}
}

// CreateRoute validates and creates a new Route
func (s *RouteUsecase) CreateRoute(ctx context.Context, request *model.WriteRouteRequest) error {
	route := mapper.ToRouteEntity(request)

	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.RouteRepository.Create(tx, route)
	})
}

// GetAllRoutes retrieves all routes
func (s *RouteUsecase) GetAllRoutes(ctx context.Context) ([]*model.ReadRouteResponse, error) {
	routes := []*entities.Route{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		routes, err = s.RouteRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all routes: %w", err)
	}

	return mapper.ToRoutesModel(routes), nil
}

// GetRouteByID retrieves a Route by its ID
func (s *RouteUsecase) GetRouteByID(ctx context.Context, id uint) (*model.ReadRouteResponse, error) {
	route := new(entities.Route)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		route, err = s.RouteRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get route by ID: %w", err)
	}

	if route == nil {
		return nil, errors.New("route not found")
	}

	return mapper.ToRouteModel(route), nil
}

// UpdateRoute updates an existing Route
func (s *RouteUsecase) UpdateRoute(ctx context.Context, id uint, request *model.WriteRouteRequest) error {
	route := mapper.ToRouteEntity(request)
	route.ID = id

	if route.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}

	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.RouteRepository.Update(tx, route)
	})
}

// DeleteRoute deletes a Route by its ID
func (s *RouteUsecase) DeleteRoute(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		route, err := s.RouteRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if route == nil {
			return errors.New("route not found")
		}
		return s.RouteRepository.Delete(tx, route)
	})
}

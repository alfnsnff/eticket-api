package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
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

func (r *RouteUsecase) CreateRoute(ctx context.Context, request *model.WriteRouteRequest) error {
	route := mapper.RouteMapper.FromWrite(request)

	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		return r.RouteRepository.Create(tx, route)
	})
}

func (r *RouteUsecase) GetAllRoutes(ctx context.Context) ([]*model.ReadRouteResponse, error) {
	routes := []*entity.Route{}

	err := tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		var err error
		routes, err = r.RouteRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all routes: %w", err)
	}

	return mapper.RouteMapper.ToModels(routes), nil
}

func (r *RouteUsecase) GetRouteByID(ctx context.Context, id uint) (*model.ReadRouteResponse, error) {
	route := new(entity.Route)

	err := tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		var err error
		route, err = r.RouteRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get route by ID: %w", err)
	}

	if route == nil {
		return nil, errors.New("route not found")
	}

	return mapper.RouteMapper.ToModel(route), nil
}

func (r *RouteUsecase) UpdateRoute(ctx context.Context, id uint, request *model.UpdateRouteRequest) error {
	route := mapper.RouteMapper.FromUpdate(request)
	route.ID = id

	if route.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}

	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		return r.RouteRepository.Update(tx, route)
	})
}

func (r *RouteUsecase) DeleteRoute(ctx context.Context, id uint) error {

	return tx.Execute(ctx, r.DB, func(tx *gorm.DB) error {
		route, err := r.RouteRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if route == nil {
			return errors.New("route not found")
		}
		return r.RouteRepository.Delete(tx, route)
	})

}

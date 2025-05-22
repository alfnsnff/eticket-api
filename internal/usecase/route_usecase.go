package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"eticket-api/pkg/utils/tx"
	"fmt"

	"gorm.io/gorm"
)

type RouteUsecase struct {
	Tx              *tx.TxManager
	RouteRepository *repository.RouteRepository
}

func NewRouteUsecase(
	tx *tx.TxManager,
	route_repository *repository.RouteRepository,
) *RouteUsecase {
	return &RouteUsecase{
		Tx:              tx,
		RouteRepository: route_repository,
	}
}

func (r *RouteUsecase) CreateRoute(ctx context.Context, request *model.WriteRouteRequest) error {
	route := mapper.RouteMapper.FromWrite(request)

	if route.DepartureHarborID == 0 || route.ArrivalHarborID == 0 {
		return fmt.Errorf("harbor ID cannot be empty")
	}

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return r.RouteRepository.Create(tx, route)
	})
}

func (r *RouteUsecase) GetAllRoutes(ctx context.Context, limit, offset int) ([]*model.ReadRouteResponse, int, error) {
	routes := []*entity.Route{}
	var total int64
	err := r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = r.RouteRepository.Count(tx)
		if err != nil {
			return err
		}
		routes, err = r.RouteRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all routes: %w", err)
	}

	return mapper.RouteMapper.ToModels(routes), int(total), nil
}

func (r *RouteUsecase) GetRouteByID(ctx context.Context, id uint) (*model.ReadRouteResponse, error) {
	route := new(entity.Route)

	err := r.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return r.RouteRepository.Update(tx, route)
	})
}

func (r *RouteUsecase) DeleteRoute(ctx context.Context, id uint) error {

	return r.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

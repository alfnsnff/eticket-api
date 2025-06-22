package route

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type RouteUsecase struct {
	DB              *gorm.DB
	RouteRepository RouteRepository
}

func NewRouteUsecase(
	db *gorm.DB,
	routeRepository RouteRepository,
) *RouteUsecase {
	return &RouteUsecase{
		DB:              db,
		RouteRepository: routeRepository,
	}
}

func (r *RouteUsecase) CreateRoute(ctx context.Context, request *model.WriteRouteRequest) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		} else {
			tx.Rollback()
		}
	}()

	route := &entity.Route{
		DepartureHarborID: request.DepartureHarborID,
		ArrivalHarborID:   request.ArrivalHarborID,
	}

	if err := r.RouteRepository.Create(tx, route); err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *RouteUsecase) GetAllRoutes(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadRouteResponse, int, error) {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		} else {
			tx.Rollback()
		}
	}()

	total, err := r.RouteRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count routes: %w", err)
	}

	routes, err := r.RouteRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all routes: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit: %w", err)
	}

	return ToReadRouteResponses(routes), int(total), nil
}

func (r *RouteUsecase) GetRouteByID(ctx context.Context, id uint) (*model.ReadRouteResponse, error) {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		} else {
			tx.Rollback()
		}
	}()

	route, err := r.RouteRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get route by ID: %w", err)
	}
	if route == nil {
		return nil, errors.New("route not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return ToReadRouteResponse(route), nil
}

func (r *RouteUsecase) UpdateRoute(ctx context.Context, request *model.UpdateRouteRequest) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	route, err := r.RouteRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find route: %w", err)
	}
	if route == nil {
		return errors.New("route not found")
	}

	route.DepartureHarborID = request.DepartureHarborID
	route.ArrivalHarborID = request.ArrivalHarborID

	if err := r.RouteRepository.Update(tx, route); err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *RouteUsecase) DeleteRoute(ctx context.Context, id uint) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		} else {
			tx.Rollback()
		}
	}()

	route, err := r.RouteRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get route: %w", err)
	}
	if route == nil {
		return errors.New("route not found")
	}

	if err := r.RouteRepository.Delete(tx, route); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

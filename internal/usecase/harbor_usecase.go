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

type HarborUsecase struct {
	DB               *gorm.DB
	HarborRepository *repository.HarborRepository
}

func NewHarborUsecase(db *gorm.DB, harbor_repository *repository.HarborRepository) *HarborUsecase {
	return &HarborUsecase{DB: db, HarborRepository: harbor_repository}
}

// Createharbor validates and creates a new harbor
func (s *HarborUsecase) CreateHarbor(ctx context.Context, request *model.WriteHarborRequest) error {
	harbor := mapper.ToHarborEntity(request)

	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.HarborRepository.Create(tx, harbor)
	})
}

// GetAllharbores retrieves all harbors
func (s *HarborUsecase) GetAllHarbors(ctx context.Context) ([]*model.ReadHarborResponse, error) {

	harbors := []*entities.Harbor{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		harbors, err = s.HarborRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ToHarborsModel(harbors), nil
}

// GetharborByID retrieves a harbor by its ID
func (s *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {

	harbor := new(entities.Harbor)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		harbor, err = s.HarborRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if harbor == nil {
		return nil, errors.New("harbor not found")
	}
	return mapper.ToHarborModel(harbor), nil
}

// Updateharbor updates an existing harbor
func (s *HarborUsecase) UpdateHarbor(ctx context.Context, id uint, request *model.WriteHarborRequest) error {
	harbor := mapper.ToHarborEntity(request)

	harbor.ID = id

	if harbor.ID == 0 {
		return fmt.Errorf("harbor ID cannot be zero")
	}
	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.HarborRepository.Update(tx, harbor)
	})
}

// Deleteharbor deletes a harbor by its ID
func (s *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		harbor, err := s.HarborRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if harbor == nil {
			return errors.New("route not found")
		}
		return s.HarborRepository.Delete(tx, harbor)
	})
}

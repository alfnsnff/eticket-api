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

type HarborUsecase struct {
	DB               *gorm.DB
	HarborRepository *repository.HarborRepository
}

func NewHarborUsecase(db *gorm.DB, harborRepository *repository.HarborRepository) *HarborUsecase {
	return &HarborUsecase{DB: db, HarborRepository: harborRepository}
}

// Createharbor validates and creates a new harbor
func (s *HarborUsecase) CreateHarbor(ctx context.Context, harbor *entities.Harbor) error {
	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.HarborRepository.Create(txDB, harbor)
	})
}

// GetAllharbores retrieves all harbors
func (s *HarborUsecase) GetAllHarbors(ctx context.Context) ([]*entities.Harbor, error) {

	var harbors []*entities.Harbor

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		harbors, err = s.HarborRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return harbors, nil
}

// GetharborByID retrieves a harbor by its ID
func (s *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*entities.Harbor, error) {

	var harbor *entities.Harbor

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		harbor, err = s.HarborRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if harbor == nil {
		return nil, errors.New("harbor not found")
	}
	return harbor, nil
}

// Updateharbor updates an existing harbor
func (s *HarborUsecase) UpdateHarbor(ctx context.Context, id uint, harbor *entities.Harbor) error {
	harbor.ID = id

	if harbor.ID == 0 {
		return fmt.Errorf("harbor ID cannot be zero")
	}
	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.HarborRepository.Update(txDB, harbor)
	})
}

// Deleteharbor deletes a harbor by its ID
func (s *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		harbor, err := s.HarborRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if harbor == nil {
			return errors.New("route not found")
		}
		return s.HarborRepository.Delete(txDB, id)
	})
}

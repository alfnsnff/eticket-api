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

type HarborUsecase struct {
	DB               *gorm.DB
	HarborRepository *repository.HarborRepository
}

func NewHarborUsecase(db *gorm.DB, harbor_repository *repository.HarborRepository) *HarborUsecase {
	return &HarborUsecase{DB: db, HarborRepository: harbor_repository}
}

func (h *HarborUsecase) CreateHarbor(ctx context.Context, request *model.WriteHarborRequest) error {
	harbor := mapper.HarborMapper.FromWrite(request)

	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, h.DB, func(tx *gorm.DB) error {
		return h.HarborRepository.Create(tx, harbor)
	})
}

func (h *HarborUsecase) GetAllHarbors(ctx context.Context) ([]*model.ReadHarborResponse, error) {
	harbors := []*entity.Harbor{}

	err := tx.Execute(ctx, h.DB, func(tx *gorm.DB) error {
		var err error
		harbors, err = h.HarborRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.HarborMapper.ToModels(harbors), nil
}

func (h *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {
	harbor := new(entity.Harbor)

	err := tx.Execute(ctx, h.DB, func(tx *gorm.DB) error {
		var err error
		harbor, err = h.HarborRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if harbor == nil {
		return nil, errors.New("harbor not found")
	}

	return mapper.HarborMapper.ToModel(harbor), nil
}

func (h *HarborUsecase) UpdateHarbor(ctx context.Context, id uint, request *model.UpdateHarborRequest) error {
	harbor := mapper.HarborMapper.FromUpdate(request)
	harbor.ID = id

	if harbor.ID == 0 {
		return fmt.Errorf("harbor ID cannot be zero")
	}
	if harbor.Name == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return tx.Execute(ctx, h.DB, func(tx *gorm.DB) error {
		return h.HarborRepository.Update(tx, harbor)
	})
}

func (h *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {

	return tx.Execute(ctx, h.DB, func(tx *gorm.DB) error {
		harbor, err := h.HarborRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if harbor == nil {
			return errors.New("route not found")
		}
		return h.HarborRepository.Delete(tx, harbor)
	})

}

package harbor

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type HarborUsecase struct {
	Tx               *tx.TxManager
	HarborRepository *repository.HarborRepository
}

func NewHarborUsecase(
	tx *tx.TxManager,
	harbor_repository *repository.HarborRepository,
) *HarborUsecase {
	return &HarborUsecase{
		Tx:               tx,
		HarborRepository: harbor_repository,
	}
}

func (h *HarborUsecase) CreateHarbor(ctx context.Context, request *model.WriteHarborRequest) error {
	harbor := mapper.HarborMapper.FromWrite(request)

	if harbor.HarborName == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return h.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return h.HarborRepository.Create(tx, harbor)
	})
}

func (h *HarborUsecase) GetAllHarbors(ctx context.Context, limit, offset int) ([]*model.ReadHarborResponse, int, error) {
	harbors := []*entity.Harbor{}
	var total int64
	err := h.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = h.HarborRepository.Count(tx)
		if err != nil {
			return err
		}
		harbors, err = h.HarborRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.HarborMapper.ToModels(harbors), int(total), nil
}

func (h *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {
	harbor := new(entity.Harbor)

	err := h.Tx.Execute(ctx, func(tx *gorm.DB) error {
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
	if harbor.HarborName == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}

	return h.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return h.HarborRepository.Update(tx, harbor)
	})
}

func (h *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {

	return h.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

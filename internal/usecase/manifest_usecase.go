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

type ManifestUsecase struct {
	Tx                 tx.TxManager
	ManifestRepository *repository.ManifestRepository
}

func NewManifestUsecase(
	tx tx.TxManager,
	manifest_repository *repository.ManifestRepository,
) *ManifestUsecase {
	return &ManifestUsecase{
		Tx:                 tx,
		ManifestRepository: manifest_repository,
	}
}

func (m *ManifestUsecase) CreateManifest(ctx context.Context, request *model.WriteManifestRequest) error {
	manifest := mapper.ManifestMapper.FromWrite(request)

	if manifest.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}

	if manifest.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	return m.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return m.ManifestRepository.Create(tx, manifest)
	})
}

func (m *ManifestUsecase) GetAllManifests(ctx context.Context, limit, offset int) ([]*model.ReadManifestResponse, int, error) {
	manifests := []*entity.Manifest{}
	var total int64
	err := m.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = m.ManifestRepository.Count(tx)
		if err != nil {
			return err
		}
		manifests, err = m.ManifestRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ManifestMapper.ToModels(manifests), int(total), nil
}

func (m *ManifestUsecase) GetManifestByID(ctx context.Context, id uint) (*model.ReadManifestResponse, error) {
	manifest := new(entity.Manifest)

	err := m.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		manifest, err = m.ManifestRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if manifest == nil {
		return nil, errors.New("ship class not found")
	}

	return mapper.ManifestMapper.ToModel(manifest), nil
}

func (m *ManifestUsecase) UpdateManifest(ctx context.Context, id uint, request *model.UpdateManifestRequest) error {
	manifest := mapper.ManifestMapper.FromUpdate(request)
	manifest.ID = id

	if manifest.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}
	if manifest.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if manifest.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}
	return m.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return m.ManifestRepository.Update(tx, manifest)
	})
}

func (m *ManifestUsecase) DeleteManifest(ctx context.Context, id uint) error {

	return m.Tx.Execute(ctx, func(tx *gorm.DB) error {
		shipClass, err := m.ManifestRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if shipClass == nil {
			return errors.New("route not found")
		}
		return m.ManifestRepository.Delete(tx, shipClass)
	})
}

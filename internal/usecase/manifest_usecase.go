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

type ManifestUsecase struct {
	DB                 *gorm.DB
	ManifestRepository *repository.ManifestRepository
}

func NewManifestUsecase(db *gorm.DB, manifest_repository *repository.ManifestRepository) *ManifestUsecase {
	return &ManifestUsecase{DB: db, ManifestRepository: manifest_repository}
}

func (m *ManifestUsecase) CreateManifest(ctx context.Context, request *model.WriteManifestRequest) error {
	manifest := mapper.ManifestMapper.FromWrite(request)

	if manifest.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}

	if manifest.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	return tx.Execute(ctx, m.DB, func(tx *gorm.DB) error {
		return m.ManifestRepository.Create(tx, manifest)
	})
}

func (m *ManifestUsecase) GetAllManifests(ctx context.Context) ([]*model.ReadManifestResponse, error) {
	manifests := []*entity.Manifest{}

	err := tx.Execute(ctx, m.DB, func(tx *gorm.DB) error {
		var err error
		manifests, err = m.ManifestRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ManifestMapper.ToModels(manifests), nil
}

func (m *ManifestUsecase) GetManifestByID(ctx context.Context, id uint) (*model.ReadManifestResponse, error) {
	manifest := new(entity.Manifest)

	err := tx.Execute(ctx, m.DB, func(tx *gorm.DB) error {
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
	return tx.Execute(ctx, m.DB, func(tx *gorm.DB) error {
		return m.ManifestRepository.Update(tx, manifest)
	})
}

func (m *ManifestUsecase) DeleteManifest(ctx context.Context, id uint) error {

	return tx.Execute(ctx, m.DB, func(tx *gorm.DB) error {
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

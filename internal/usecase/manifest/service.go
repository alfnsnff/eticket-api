package manifest

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type ManifestUsecase struct {
	DB                 *gorm.DB
	ManifestRepository *repository.ManifestRepository
}

func NewManifestUsecase(
	db *gorm.DB,
	manifestRepository *repository.ManifestRepository,
) *ManifestUsecase {
	return &ManifestUsecase{
		DB:                 db,
		ManifestRepository: manifestRepository,
	}
}

func (m *ManifestUsecase) CreateManifest(ctx context.Context, request *model.WriteManifestRequest) error {
	tx := m.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	manifest := &entity.Manifest{
		ShipID:   request.ShipID,
		ClassID:  request.ClassID,
		Capacity: request.Capacity,
	}

	if manifest.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if manifest.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	if err := m.ManifestRepository.Create(tx, manifest); err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *ManifestUsecase) GetAllManifests(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadManifestResponse, int, error) {
	tx := m.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := m.ManifestRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count manifests: %w", err)
	}

	manifests, err := m.ManifestRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get manifests: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ManifestMapper.ToModels(manifests), int(total), nil
}

func (m *ManifestUsecase) GetManifestByID(ctx context.Context, id uint) (*model.ReadManifestResponse, error) {
	tx := m.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	manifest, err := m.ManifestRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}
	if manifest == nil {
		return nil, errors.New("ship class not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ManifestMapper.ToModel(manifest), nil
}

func (m *ManifestUsecase) UpdateManifest(ctx context.Context, request *model.UpdateManifestRequest) error {
	tx := m.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if request.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}

	// Fetch existing allocation
	manifest, err := m.ManifestRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find manifest: %w", err)
	}
	if manifest == nil {
		return errors.New("manifest not found")
	}

	manifest.ClassID = request.ClassID
	manifest.ShipID = request.ShipID
	manifest.Capacity = request.Capacity

	if err := m.ManifestRepository.Update(tx, manifest); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *ManifestUsecase) DeleteManifest(ctx context.Context, id uint) error {
	tx := m.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	manifest, err := m.ManifestRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %w", err)
	}
	if manifest == nil {
		return errors.New("manifest not found")
	}

	if err := m.ManifestRepository.Delete(tx, manifest); err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

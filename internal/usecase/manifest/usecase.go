package manifest

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type ManifestUsecase struct {
	DB                 *gorm.DB
	ManifestRepository ManifestRepository
}

func NewManifestUsecase(
	db *gorm.DB,
	manifestRepository ManifestRepository,
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

	manifest := &domain.Manifest{
		ShipID:   request.ShipID,
		ClassID:  request.ClassID,
		Capacity: request.Capacity,
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

	responses := make([]*model.ReadManifestResponse, len(manifests))
	for i, manifest := range manifests {
		responses[i] = ManifestToResponse(manifest)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
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

	return ManifestToResponse(manifest), nil
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

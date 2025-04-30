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

type ManifestUsecase struct {
	DB                 *gorm.DB
	ManifestRepository *repository.ManifestRepository
}

func NewManifestUsecase(db *gorm.DB, manifest_repository *repository.ManifestRepository) *ManifestUsecase {
	return &ManifestUsecase{DB: db, ManifestRepository: manifest_repository}
}

func (s *ManifestUsecase) CreateManifest(ctx context.Context, request *model.WriteManifestRequest) error {
	manifest := mapper.ToManifestEntity(request)

	if manifest.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if manifest.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	existing := new(entities.Manifest)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		existing, err = s.ManifestRepository.GetByShipAndClass(tx, manifest.ShipID, manifest.ClassID)
		return err
	})

	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("ship class with this ship ID and class ID already exists")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ManifestRepository.Create(tx, manifest)
	})
}

// GetAllshipes retrieves all ships
func (s *ManifestUsecase) GetAllManifests(ctx context.Context) ([]*model.ReadManifestResponse, error) {
	// return s.ManifestRepository.GetAll()

	manifests := []*entities.Manifest{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		manifests, err = s.ManifestRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ToManifestsModel(manifests), nil
}

// GetshipByID retrieves a ship by its ID
func (s *ManifestUsecase) GetManifestByID(ctx context.Context, id uint) (*model.ReadManifestResponse, error) {

	manifest := new(entities.Manifest)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		manifest, err = s.ManifestRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if manifest == nil {
		return nil, errors.New("ship class not found")
	}
	return mapper.ToManifestModel(manifest), nil
}

// GetshipByID retrieves a ship by its ID
func (s *ManifestUsecase) GetManifestsByShipID(ctx context.Context, shipId uint) ([]*model.ReadManifestResponse, error) {
	// shipClasses, err := s.ManifestRepository.GetByShipID(shipId)

	manifests := []*entities.Manifest{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		manifests, err = s.ManifestRepository.GetByShipID(tx, shipId)
		return err
	})

	if err != nil {
		return nil, err
	}
	if manifests == nil {
		return nil, errors.New("ship class not found")
	}
	return mapper.ToManifestsModel(manifests), nil
}

// Updateship updates an existing ship
func (s *ManifestUsecase) UpdateManifest(ctx context.Context, id uint, request *model.WriteManifestRequest) error {
	manifest := mapper.ToManifestEntity(request)
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
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ManifestRepository.Update(tx, manifest)
	})
}

// Deleteship deletes a ship by its ID
func (s *ManifestUsecase) DeleteManifest(ctx context.Context, id uint) error {

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		shipClass, err := s.ManifestRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if shipClass == nil {
			return errors.New("route not found")
		}
		return s.ManifestRepository.Delete(tx, shipClass)
	})
}

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

type ShipClassUsecase struct {
	DB                  *gorm.DB
	ShipClassRepository *repository.ShipClassRepository
}

func NewShipClassUsecase(db *gorm.DB, shipClassRepository *repository.ShipClassRepository) *ShipClassUsecase {
	return &ShipClassUsecase{DB: db, ShipClassRepository: shipClassRepository}
}

func (s *ShipClassUsecase) CreateShipClass(ctx context.Context, shipClass *entities.ShipClass) error {
	if shipClass.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if shipClass.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	var existing *entities.ShipClass

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		existing, err = s.ShipClassRepository.GetByShipAndClass(txDB, shipClass.ShipID, shipClass.ClassID)
		return err
	})

	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("ship class with this ship ID and class ID already exists")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ShipClassRepository.Create(txDB, shipClass)
	})
}

// GetAllshipes retrieves all ships
func (s *ShipClassUsecase) GetAllShipClasses(ctx context.Context) ([]*entities.ShipClass, error) {
	// return s.ShipClassRepository.GetAll()

	var shipClasses []*entities.ShipClass

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		shipClasses, err = s.ShipClassRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return shipClasses, nil
}

// GetshipByID retrieves a ship by its ID
func (s *ShipClassUsecase) GetShipClassByID(ctx context.Context, id uint) (*entities.ShipClass, error) {

	var shipClass *entities.ShipClass

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		shipClass, err = s.ShipClassRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if shipClass == nil {
		return nil, errors.New("ship class not found")
	}
	return shipClass, nil
}

// GetshipByID retrieves a ship by its ID
func (s *ShipClassUsecase) GetShipClassByShipID(ctx context.Context, shipId uint) ([]*entities.ShipClass, error) {
	// shipClasses, err := s.ShipClassRepository.GetByShipID(shipId)

	var shipClasses []*entities.ShipClass

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		shipClasses, err = s.ShipClassRepository.GetByShipID(txDB, shipId)
		return err
	})

	if err != nil {
		return nil, err
	}
	if shipClasses == nil {
		return nil, errors.New("ship class not found")
	}
	return shipClasses, nil
}

// Updateship updates an existing ship
func (s *ShipClassUsecase) UpdateShipClass(ctx context.Context, id uint, shipClass *entities.ShipClass) error {
	shipClass.ID = id

	if shipClass.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}
	if shipClass.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if shipClass.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ShipClassRepository.Update(txDB, shipClass)
	})
}

// Deleteship deletes a ship by its ID
func (s *ShipClassUsecase) DeleteShipClass(ctx context.Context, id uint) error {

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		shipClass, err := s.ShipClassRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if shipClass == nil {
			return errors.New("route not found")
		}
		return s.ShipClassRepository.Delete(txDB, id)
	})
}

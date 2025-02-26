package usecase

import (
	"errors"
	"eticket-api/internal/domain"
	"fmt"
)

type HarborUsecase struct {
	HarborRepository domain.HarborRepositoryInterface
}

func NewHarborUsecase(harborRepository domain.HarborRepositoryInterface) HarborUsecase {
	return HarborUsecase{HarborRepository: harborRepository}
}

// Createharbor validates and creates a new harbor
func (s *HarborUsecase) CreateHarbor(harbor *domain.Harbor) error {
	if harbor.HarborName == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}
	return s.HarborRepository.Create(harbor)
}

// GetAllharbores retrieves all harbors
func (s *HarborUsecase) GetAllHarbors() ([]*domain.Harbor, error) {
	return s.HarborRepository.GetAll()
}

// GetharborByID retrieves a harbor by its ID
func (s *HarborUsecase) GetHarborByID(id uint) (*domain.Harbor, error) {
	harbor, err := s.HarborRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if harbor == nil {
		return nil, errors.New("harbor not found")
	}
	return harbor, nil
}

// Updateharbor updates an existing harbor
func (s *HarborUsecase) UpdateHarbor(harbor *domain.Harbor) error {
	if harbor.ID == 0 {
		return fmt.Errorf("harbor ID cannot be zero")
	}
	if harbor.HarborName == "" {
		return fmt.Errorf("harbor name cannot be empty")
	}
	return s.HarborRepository.Update(harbor)
}

// Deleteharbor deletes a harbor by its ID
func (s *HarborUsecase) DeleteHarbor(id uint) error {
	harbor, err := s.HarborRepository.GetByID(id)
	if err != nil {
		return err
	}
	if harbor == nil {
		return errors.New("harbor not found")
	}
	return s.HarborRepository.Delete(id)
}

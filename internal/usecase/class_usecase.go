package usecase

import (
	"errors"
	"eticket-api/internal/domain"
	"fmt"
)

type ClassUsecase struct {
	ClassRepository domain.ClassRepository
}

func NewClassUsecase(classRepository domain.ClassRepository) ClassUsecase {
	return ClassUsecase{ClassRepository: classRepository}
}

// CreateClass validates and creates a new class
func (s *ClassUsecase) CreateClass(class *domain.Class) error {
	if class.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if class.ClassName == "" {
		return fmt.Errorf("class name cannot be empty")
	}
	return s.ClassRepository.Create(class)
}

// GetAllClasses retrieves all classes
func (s *ClassUsecase) GetAllClasses() ([]*domain.Class, error) {
	return s.ClassRepository.GetAll()
}

// GetClassByID retrieves a class by its ID
func (s *ClassUsecase) GetClassByID(id uint) (*domain.Class, error) {
	class, err := s.ClassRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}
	return class, nil
}

// UpdateClass updates an existing class
func (s *ClassUsecase) UpdateClass(class *domain.Class) error {
	if class.ID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	if class.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if class.ClassName == "" {
		return fmt.Errorf("class name cannot be empty")
	}
	return s.ClassRepository.Update(class)
}

// DeleteClass deletes a class by its ID
func (s *ClassUsecase) DeleteClass(id uint) error {
	class, err := s.ClassRepository.GetByID(id)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}
	return s.ClassRepository.Delete(id)
}

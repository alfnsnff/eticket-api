package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	"fmt"
)

type ClassUsecase struct {
	ClassRepository *repository.ClassRepository
}

func NewClassUsecase(classRepository *repository.ClassRepository) ClassUsecase {
	return ClassUsecase{ClassRepository: classRepository}
}

// CreateClass validates and creates a new class
func (s *ClassUsecase) CreateClass(class *entities.Class) error {
	if class.Name == "" {
		return fmt.Errorf("class name cannot be empty")
	}
	return s.ClassRepository.Create(class)
}

// GetAllClasses retrieves all classes
func (s *ClassUsecase) GetAllClasses() ([]*entities.Class, error) {
	return s.ClassRepository.GetAll()
}

// GetClassByID retrieves a class by its ID
func (s *ClassUsecase) GetClassByID(id uint) (*entities.Class, error) {
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
func (s *ClassUsecase) UpdateClass(id uint, class *entities.Class) error {
	class.ID = id

	if class.ID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	if class.Name == "" {
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

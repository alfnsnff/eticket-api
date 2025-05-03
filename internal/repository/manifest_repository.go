package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ManifestRepository struct {
	Repository[entity.Manifest]
}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}

func (mr *ManifestRepository) GetAll(db *gorm.DB) ([]*entity.Manifest, error) {
	manifests := []*entity.Manifest{}
	result := db.Find(&manifests)
	if result.Error != nil {
		return nil, result.Error
	}
	return manifests, nil
}

func (mr *ManifestRepository) GetByID(db *gorm.DB, id uint) (*entity.Manifest, error) {
	manifest := new(entity.Manifest)
	result := db.First(&manifest, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manifest, result.Error
}

func (mr *ManifestRepository) GetByShipAndClass(ctx *gorm.DB, shipID uint, classID uint) (*entity.Manifest, error) {
	manifest := new(entity.Manifest)
	result := ctx.Where("ship_id = ? AND class_id = ?", shipID, classID).First(manifest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manifest, result.Error
}

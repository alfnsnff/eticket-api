package repository

import (
	"errors"
	"eticket-api/internal/entity"

	"gorm.io/gorm"
)

type ManifestRepository struct {
	Repository[entity.Manifest]
}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}

func (mr *ManifestRepository) Count(db *gorm.DB) (int64, error) {
	manifests := []*entity.Manifest{}
	var total int64
	result := db.Find(&manifests).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (mr *ManifestRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Manifest, error) {
	manifests := []*entity.Manifest{}
	result := db.Preload("Class").
		Preload("Ship").
		Limit(limit).Offset(offset).Find(&manifests)
	if result.Error != nil {
		return nil, result.Error
	}
	return manifests, nil
}

func (mr *ManifestRepository) GetByID(db *gorm.DB, id uint) (*entity.Manifest, error) {
	manifest := new(entity.Manifest)
	result := db.Preload("Class").
		Preload("Ship").
		First(&manifest, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manifest, result.Error
}

func (mr *ManifestRepository) GetByShipAndClass(ctx *gorm.DB, shipID uint, classID uint) (*entity.Manifest, error) {
	manifest := new(entity.Manifest)
	result := ctx.Preload("Class").
		Preload("Ship").
		Where("ship_id = ? AND class_id = ?", shipID, classID).First(manifest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manifest, result.Error
}

func (r *ManifestRepository) FindByShipID(db *gorm.DB, shipID uint) ([]*entity.Manifest, error) {
	manifests := []*entity.Manifest{}
	result := db.Preload("Class").
		Preload("Ship").
		Where("ship_id = ?", shipID).Find(&manifests)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return manifests, nil
}

package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManifestRepository struct{}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}

func (ar *ManifestRepository) Create(db *gorm.DB, manifest *entity.Manifest) error {
	result := db.Create(manifest)
	return result.Error
}

func (ar *ManifestRepository) Update(db *gorm.DB, manifest *entity.Manifest) error {
	result := db.Save(manifest)
	return result.Error
}

func (ar *ManifestRepository) Delete(db *gorm.DB, manifest *entity.Manifest) error {
	result := db.Select(clause.Associations).Delete(manifest)
	return result.Error
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

func (mr *ManifestRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Manifest, error) {
	manifests := []*entity.Manifest{}

	query := db.Preload("Class").Preload("Ship")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("ship_id ? OR class_id ILIKE ?", search, search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&manifests).Error
	return manifests, err
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

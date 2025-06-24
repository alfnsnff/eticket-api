package domain

import (
	"time"

	"gorm.io/gorm"
)

type Manifest struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	ShipID    uint      `gorm:"column:ship_id;not null;index;uniqueIndex:idx_ship_class"`
	ClassID   uint      `gorm:"column:class_id;not null;index;uniqueIndex:idx_ship_class"`
	Capacity  int       `gorm:"column:capacity;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"`
	Ship  Ship  `gorm:"foreignKey:ShipID"`
}

func (m *Manifest) TableName() string {
	return "manifest"
}

type ManifestRepository interface {
	Create(db *gorm.DB, entity *Manifest) error
	Update(db *gorm.DB, entity *Manifest) error
	Delete(db *gorm.DB, entity *Manifest) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Manifest, error)
	GetByID(db *gorm.DB, id uint) (*Manifest, error)
	GetByShipAndClass(db *gorm.DB, shipID uint, classID uint) (*Manifest, error)
	FindByShipID(db *gorm.DB, shipID uint) ([]*Manifest, error)
}

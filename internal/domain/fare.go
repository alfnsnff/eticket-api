package domain

import (
	"time"

	"gorm.io/gorm"
)

type Fare struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	RouteID     uint      `gorm:"column:route_id;not null;index;"`
	ManifestID  uint      `gorm:"column:manifest_id;not null;index;"`
	TicketPrice float32   `gorm:"column:ticket_price;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`

	Route    Route    `gorm:"foreignKey:RouteID"`
	Manifest Manifest `gorm:"foreignKey:ManifestID"`
}

func (f *Fare) TableName() string {
	return "fare"
}

type FareRepository interface {
	Create(db *gorm.DB, entity *Fare) error
	Update(db *gorm.DB, entity *Fare) error
	Delete(db *gorm.DB, entity *Fare) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Fare, error)
	GetByID(db *gorm.DB, id uint) (*Fare, error)
	GetByManifestAndRoute(db *gorm.DB, manifestID uint, routeID uint) (*Fare, error)
}

package domain

import (
	"time"

	"gorm.io/gorm"
)

type Ship struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id"`
	ShipName      string    `gorm:"column:ship_name;unique;not null"`
	Status        string    `gorm:"column:status;type:varchar(24);not null"`
	ShipType      string    `gorm:"column:ship_type;type:varchar(24);not null"`
	ShipAlias     string    `gorm:"column:ship_alias;type:varchar(8);not null"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	ImageLink     string    `gorm:"column:image_link;not null"`
	Description   string    `gorm:"column:description;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (sh *Ship) TableName() string {
	return "ship"
}

type ShipRepository interface {
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *Ship) error
	InsertBulk(db *gorm.DB, ships []*Ship) error
	Update(db *gorm.DB, entity *Ship) error
	UpdateBulk(db *gorm.DB, ships []*Ship) error
	Delete(db *gorm.DB, entity *Ship) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Ship, error)
	FindByID(db *gorm.DB, id uint) (*Ship, error)
}

package domain

import (
	"time"

	"gorm.io/gorm"
)

type Harbor struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	HarborName    string    `gorm:"column:harbor_name;type:varchar(24);unique;;not null"`
	Status        string    `gorm:"column:harbor_status;idtype:varchar(24);not null"`
	HarborAlias   string    `gorm:"column:harbor_alias;type:varchar(8);"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (h *Harbor) TableName() string {
	return "harbor"
}

type HarborRepository interface {
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *Harbor) error
	InsertBulk(db *gorm.DB, harbors []*Harbor) error
	Update(db *gorm.DB, entity *Harbor) error
	UpdateBulk(db *gorm.DB, harbors []*Harbor) error
	Delete(db *gorm.DB, entity *Harbor) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Harbor, error)
	FindByID(db *gorm.DB, id uint) (*Harbor, error)
}

package domain

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ClassName  string    `gorm:"column:class_name;type:varchar(24);not null"`
	Type       string    `gorm:"column:type;type:varchar(24);not null"`
	ClassAlias *string   `gorm:"column:class_alias;type:varchar(8);"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

func (c *Class) TableName() string {
	return "class"
}

type ClassRepository interface {
	Create(db *gorm.DB, entity *Class) error
	Update(db *gorm.DB, entity *Class) error
	Delete(db *gorm.DB, entity *Class) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Class, error)
	GetByID(db *gorm.DB, id uint) (*Class, error)
}

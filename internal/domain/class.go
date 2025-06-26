package domain

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ClassName  string    `gorm:"column:class_name;unique;type:varchar(24);not null"`
	Type       string    `gorm:"column:type;type:varchar(24);not null"`
	ClassAlias string    `gorm:"column:class_alias;type:varchar(8);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

func (c *Class) TableName() string {
	return "class"
}

type ClassRepository interface {
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *Class) error
	InsertBulk(db *gorm.DB, classes []*Class) error
	Update(db *gorm.DB, entity *Class) error
	UpdateBulk(db *gorm.DB, classes []*Class) error
	Delete(db *gorm.DB, entity *Class) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Class, error)
	FindByID(db *gorm.DB, id uint) (*Class, error)
}

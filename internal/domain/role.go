package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint      `gorm:"primaryKey"`
	RoleName    string    `gorm:"column:role_name;unique;not null"` // e.g., "admin", "editor"
	Description string    `gorm:"column:description;type:varchar(128);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (r *Role) TableName() string {
	return "role"
}

type RoleRepository interface {
	Create(db *gorm.DB, entity *Role) error
	Update(db *gorm.DB, entity *Role) error
	Delete(db *gorm.DB, entity *Role) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Role, error)
	GetByID(db *gorm.DB, id uint) (*Role, error)
}

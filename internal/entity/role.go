package entity

import "time"

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

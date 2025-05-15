package entity

import "time"

type UserRole struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"column:user_id;not null"`
	RoleID    uint      `gorm:"column:role_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (ur *UserRole) TableName() string {
	return "user_role"
}

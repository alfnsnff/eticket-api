package entity

import "time"

type Class struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	ClassName string    `gorm:"column:class_name;type:varchar(24);not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (c *Class) TableName() string {
	return "class"
}

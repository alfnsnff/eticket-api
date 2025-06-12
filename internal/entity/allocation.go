package entity

import (
	"time"
)

type Allocation struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ScheduleID uint      `gorm:"column:schedule_id;not null"`
	ClassID    uint      `gorm:"column:class_id;not null"`
	Quota      int       `gorm:"column:quota;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"`
}

func (a *Allocation) TableName() string {
	return "allocation"
}

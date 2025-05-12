package entity

import "time"

type Manifest struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	ShipID    uint      `gorm:"column:ship_id;not null;index;" json:"ship_id"`
	ClassID   uint      `gorm:"column:class_id;not null;index;" json:"class_id"`
	Capacity  int       `gorm:"column:capacity;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (m *Manifest) TableName() string {
	return "manifest"
}

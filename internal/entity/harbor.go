package entity

import "time"

type Harbor struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	HarborName    string    `gorm:"column:harbor_name;type:varchar(24);not null"`
	Status        string    `gorm:"column:harbor_status;idtype:varchar(24);not null"`
	HarborAlias   *string   `gorm:"column:harbor_alias;type:varchar(8);"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (h *Harbor) TableName() string {
	return "harbor"
}

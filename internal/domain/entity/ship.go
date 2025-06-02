package entity

import "time"

type Ship struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id"`
	ShipName      string    `gorm:"column:ship_name;not null"`
	Status        string    `gorm:"column:status;type:varchar(24);not null"`
	ShipType      string    `gorm:"column:ship_type;type:varchar(24);not null"`
	ShipAlias     *string   `gorm:"column:ship_alias;type:varchar(8);"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	ImageLink     string    `gorm:"column:image_link;not null"`
	Description   string    `gorm:"column:description;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (sh *Ship) TableName() string {
	return "ship"
}

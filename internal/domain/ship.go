package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type Ship struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id"`
	ShipName      string    `gorm:"column:ship_name;unique;not null"`
	Status        string    `gorm:"column:status;type:varchar(24);not null"`
	ShipType      string    `gorm:"column:ship_type;type:varchar(24);not null"`
	ShipAlias     string    `gorm:"column:ship_alias;type:varchar(8);not null"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	ImageLink     string    `gorm:"column:image_link;not null"`
	Description   string    `gorm:"column:description;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (sh *Ship) TableName() string {
	return "ship"
}

type ShipRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Ship) error
	InsertBulk(ctx context.Context, conn gotann.Connection, ships []*Ship) error
	Update(ctx context.Context, conn gotann.Connection, entity *Ship) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, ships []*Ship) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Ship) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Ship, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Ship, error)
}

package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type Harbor struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	HarborName    string    `gorm:"column:harbor_name;type:varchar(24);unique;;not null"`
	Status        string    `gorm:"column:harbor_status;idtype:varchar(24);not null"`
	HarborAlias   string    `gorm:"column:harbor_alias;type:varchar(8);"`
	YearOperation string    `gorm:"column:year_operation;type:varchar(24);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (h *Harbor) TableName() string {
	return "harbor"
}

type HarborRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Harbor) error
	InsertBulk(ctx context.Context, conn gotann.Connection, harbors []*Harbor) error
	Update(ctx context.Context, conn gotann.Connection, entity *Harbor) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, harbors []*Harbor) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Harbor) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Harbor, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Harbor, error)
}

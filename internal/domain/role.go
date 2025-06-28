package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
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
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Role) error
	InsertBulk(ctx context.Context, conn gotann.Connection, roles []*Role) error
	Update(ctx context.Context, conn gotann.Connection, entity *Role) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, roles []*Role) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Role) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Role, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Role, error)
}

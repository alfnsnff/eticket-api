package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type Class struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ClassName  string    `gorm:"column:class_name;unique;type:varchar(24);not null"`
	Type       string    `gorm:"column:type;type:varchar(24);not null"`
	ClassAlias string    `gorm:"column:class_alias;type:varchar(8);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

func (c *Class) TableName() string {
	return "class"
}

type ClassRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Class) error
	InsertBulk(ctx context.Context, conn gotann.Connection, classes []*Class) error
	Update(ctx context.Context, conn gotann.Connection, entity *Class) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, classes []*Class) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Class) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Class, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Class, error)
}

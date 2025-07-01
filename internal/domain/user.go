package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	RoleID    uint      `gorm:"column:role_id;not null"`
	Username  string    `gorm:"column:username;type:varchar(16);unique;not null"` // Used as Casbin subject
	Email     string    `gorm:"column:email;unique;not null"`
	Password  string    `gorm:"column:password;not null"` // Store bcrypt hash
	FullName  string    `gorm:"column:full_name;type:varchar(32);not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	Role Role `gorm:"foreignKey:RoleID"`
}

func (u *User) TableName() string {
	return "user"
}

type UserRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *User) error
	InsertBulk(ctx context.Context, conn gotann.Connection, users []*User) error
	Update(ctx context.Context, conn gotann.Connection, entity *User) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, users []*User) error
	Delete(ctx context.Context, conn gotann.Connection, entity *User) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*User, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*User, error)
	FindByEmail(ctx context.Context, conn gotann.Connection, email string) (*User, error)
	FindByUsername(ctx context.Context, conn gotann.Connection, username string) (*User, error)
}

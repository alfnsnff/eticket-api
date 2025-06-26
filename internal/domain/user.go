package domain

import (
	"time"

	"gorm.io/gorm"
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
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *User) error
	InsertBulk(db *gorm.DB, users []*User) error
	Update(db *gorm.DB, entity *User) error
	UpdateBulk(db *gorm.DB, users []*User) error
	Delete(db *gorm.DB, entity *User) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*User, error)
	FindByID(db *gorm.DB, id uint) (*User, error)
	FindByEmail(db *gorm.DB, email string) (*User, error)
	FindByUsername(db *gorm.DB, username string) (*User, error)
}

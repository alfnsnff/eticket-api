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
	Create(db *gorm.DB, entity *User) error
	Update(db *gorm.DB, entity *User) error
	Delete(db *gorm.DB, entity *User) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*User, error)
	GetByID(db *gorm.DB, id uint) (*User, error)
	GetByUsername(db *gorm.DB, username string) (*User, error)
	GetByEmail(db *gorm.DB, email string) (*User, error)
	UpdatePassword(db *gorm.DB, userID uint, password string) error
}

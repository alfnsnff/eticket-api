package entity

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"` // Used as Casbin subject
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"` // Store bcrypt hash
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "user"
}

package entity

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	RoleID uint
}

func (ur *UserRole) TableName() string {
	return "user_role"
}

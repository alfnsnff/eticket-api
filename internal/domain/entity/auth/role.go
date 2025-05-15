package entity

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"` // e.g., "admin", "editor"
	DisplayName string
	Description string
}

func (r *Role) TableName() string {
	return "role"
}

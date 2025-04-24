package entities

import "time"

type ShipClass struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShipID    uint      `gorm:"not null" json:"ship_id"`
	ClassID   uint      `gorm:"not null" json:"class_id"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Ship  Ship  `gorm:"foreignKey:ShipID" json:"ship"`
	Class Class `gorm:"foreignKey:ClassID" json:"class"`
}

type ShipClassRepositoryInterface interface {
	Create(shipClass *ShipClass) error
	GetAll() ([]*ShipClass, error)
	GetByShipAndClass(shipID, classID uint) (*ShipClass, error)
	GetByID(id uint) (*ShipClass, error) // Add this method
	GetByIDs(ids []uint) ([]*ShipClass, error)
	GetByShipID(shipID uint) ([]*ShipClass, error)
	Update(shipClass *ShipClass) error
	Delete(id uint) error
}

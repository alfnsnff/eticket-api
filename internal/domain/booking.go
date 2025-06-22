package domain

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	OrderID         *string   `gorm:"column:order_id;type:varchar(64);;uniqueIndex"`
	ScheduleID      uint      `gorm:"column:schedule_id;not null;index;"`
	IDType          string    `gorm:"column:id_type;type:varchar(24);not null"`
	IDNumber        string    `gorm:"column:id_number;type:varchar(24);not null"`
	CustomerName    string    `gorm:"column:customer_name;type:varchar(32);not null"`
	CustomerAge     int       `gorm:"column:customer_age;not null"`
	CustomerGender  string    `gorm:"column:customer_gender;type:varchar(24);not null"`
	PhoneNumber     string    `gorm:"column:phone_number;type:varchar(14);not null"`
	Email           string    `gorm:"column:email;not null"`
	ReferenceNumber *string   `gorm:"column:reference_number;"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`

	Tickets  []Ticket `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID"`
}

func (b *Booking) TableName() string {
	return "booking"
}

type BookingRepository interface {
	Create(db *gorm.DB, entity *Booking) error
	Update(db *gorm.DB, entity *Booking) error
	Delete(db *gorm.DB, entity *Booking) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Booking, error)
	GetByID(db *gorm.DB, id uint) (*Booking, error)
	GetByOrderID(db *gorm.DB, id string) (*Booking, error)
	PaidConfirm(db *gorm.DB, id uint) error
	UpdateReferenceNumber(tx *gorm.DB, bookingID uint, reference *string) error
}

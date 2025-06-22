package entity

import (
	"time"
)

type Ticket struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	ScheduleID      uint      `gorm:"column:schedule_id;not null;index;"`
	ClassID         uint      `gorm:"column:class_id;not null;index;"`
	BookingID       *uint     `gorm:"column:booking_id;index"`
	ClaimSessionID  *uint     `gorm:"column:claim_session_id;index"`
	PassengerName   *string   `gorm:"column:passenger_name;type:varchar(32)"`
	PassengerAge    *int      `gorm:"column:passenger_age;"`
	PassengerGender *string   `gorm:"column:passenger_gender;type:varchar(24);"`
	Address         *string   `gorm:"column:address;type:varchar(32)"`
	IDType          *string   `gorm:"column:id_type;type:varchar(24)"`
	IDNumber        *string   `gorm:"column:id_number;type:varchar(24)"`
	Type            string    `gorm:"column:type;type:varchar(20);not null"` // "passenger" or "vehicle"
	SeatNumber      *string   `gorm:"column:seat_number;type:varchar(24)"`
	LicensePlate    *string   `gorm:"column:license_plate;type:varchar(24)"`
	Price           float32   `gorm:"column:price;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`

	Class    Class    `gorm:"foreignKey:ClassID"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID"`
	Booking  Booking  `gorm:"foreignKey:BookingID"`
}

func (t *Ticket) TableName() string {
	return "ticket"
}

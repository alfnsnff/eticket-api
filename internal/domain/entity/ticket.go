package entity

import "time"

type Ticket struct {
	ID             uint       `gorm:"column:id;primaryKey"`
	ScheduleID     uint       `gorm:"column:schedule_id;not null;index;"`
	ClassID        uint       `gorm:"column:class_id;not null;index;"`
	BookingID      *uint      `gorm:"column:booking_id;index"`
	ClaimSessionID *uint      `gorm:"column:claim_session_id;index"`
	Status         string     `gorm:"column:status;type:varchar(24);not null"`
	Price          float32    `gorm:"column:price;not null"`
	PassengerName  *string    `gorm:"column:passenger_name;type:varchar(32)"`
	PassengerAge   *int       `gorm:"column:passenger_age;"`
	Address        *string    `gorm:"column:address;type:varchar(32)"`
	IDType         *string    `gorm:"column:id_type;type:varchar(24)"`
	IDNumber       *string    `gorm:"column:id_number;type:varchar(24)"`
	SeatNumber     *string    `gorm:"column:seat_number;type:varchar(8)"`
	EntriesAt      *time.Time `gorm:"column:entries_at"`
	BookedAt       *time.Time `gorm:"column:booked_at"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"`
}

func (t *Ticket) TableName() string {
	return "ticket"
}

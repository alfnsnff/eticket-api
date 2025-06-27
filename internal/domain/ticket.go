package domain

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	BookingID       *uint     `gorm:"column:booking_id;not null;index;"`
	TicketCode      string    `gorm:"column:ticket_code;type:varchar(64);not null;uniqueIndex"` // Unique ticket code
	ScheduleID      uint      `gorm:"column:schedule_id;not null;index;"`
	ClassID         uint      `gorm:"column:class_id;not null;index;"`
	PassengerName   string    `gorm:"column:passenger_name;type:varchar(32)"`
	PassengerAge    int       `gorm:"column:passenger_age;"`
	Address         string    `gorm:"column:address;type:varchar(32)"`
	PassengerGender *string   `gorm:"column:passenger_gender;type:varchar(24);"`
	IDType          *string   `gorm:"column:id_type;type:varchar(24)"`
	IDNumber        *string   `gorm:"column:id_number;type:varchar(24)"`
	SeatNumber      *string   `gorm:"column:seat_number;type:varchar(24)"`
	LicensePlate    *string   `gorm:"column:license_plate;type:varchar(24)"`
	Type            string    `gorm:"column:type;type:varchar(20);not null"` // "passenger" or "vehicle"
	Price           float64   `gorm:"column:price;not null"`
	IsCheckedIn     bool      `gorm:"column:is_checked_in;not null;default:false"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`

	Class    Class    `gorm:"foreignKey:ClassID"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID"`
	Booking  Booking  `gorm:"foreignKey:BookingID"`
}

func (t *Ticket) TableName() string {
	return "ticket"
}

type TicketRepository interface {
	Count(db *gorm.DB) (int64, error)
	CountByScheduleID(db *gorm.DB, scheduleID uint) (int64, error)
	CountByScheduleIDAndClassIDWithStatus(db *gorm.DB, scheduleID uint, classID uint) (int64, error)
	Insert(db *gorm.DB, entity *Ticket) error
	InsertBulk(db *gorm.DB, tickets []*Ticket) error
	Update(db *gorm.DB, entity *Ticket) error
	UpdateBulk(db *gorm.DB, tickets []*Ticket) error
	Delete(db *gorm.DB, entity *Ticket) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Ticket, error)
	FindByID(db *gorm.DB, id uint) (*Ticket, error)
	FindByIDs(db *gorm.DB, ids []uint) ([]*Ticket, error)
	FindByBookingID(db *gorm.DB, bookingID uint) ([]*Ticket, error)
	FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*Ticket, error)
	FindByClaimSessionID(db *gorm.DB, sessionID uint) ([]*Ticket, error)
	CheckIn(db *gorm.DB, id uint) error
}

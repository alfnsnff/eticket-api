package domain

import (
	"time"

	"gorm.io/gorm"
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
	Create(db *gorm.DB, entity *Ticket) error
	Update(db *gorm.DB, entity *Ticket) error
	Delete(db *gorm.DB, entity *Ticket) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Ticket, error)
	GetByScheduleID(db *gorm.DB, id, limit, offset int, sort, search string) ([]*Ticket, error)
	GetByID(db *gorm.DB, id uint) (*Ticket, error)
	GetByBookingID(db *gorm.DB, id uint) ([]*Ticket, error)
	CountByScheduleClassAndStatuses(db *gorm.DB, scheduleID uint, classID uint) (int64, error)
	CreateBulk(db *gorm.DB, tickets []*Ticket) error
	UpdateBulk(db *gorm.DB, tickets []*Ticket) error
	FindManyByIDs(db *gorm.DB, ids []uint) ([]*Ticket, error)
	FindManyBySessionID(db *gorm.DB, sessionID uint) ([]*Ticket, error)
	CancelManyBySessionID(db *gorm.DB, sessionID uint) error
	CheckIn(db *gorm.DB, id uint) error
}

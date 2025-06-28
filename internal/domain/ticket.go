package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
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
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	CountByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) (int64, error)
	CountByScheduleIDAndClassIDWithStatus(ctx context.Context, conn gotann.Connection, scheduleID uint, classID uint) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Ticket) error
	InsertBulk(ctx context.Context, conn gotann.Connection, tickets []*Ticket) error
	Update(ctx context.Context, conn gotann.Connection, entity *Ticket) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, tickets []*Ticket) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Ticket) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Ticket, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Ticket, error)
	FindByIDs(ctx context.Context, conn gotann.Connection, ids []uint) ([]*Ticket, error)
	FindByBookingID(ctx context.Context, conn gotann.Connection, bookingID uint) ([]*Ticket, error)
	FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*Ticket, error)
	FindByClaimSessionID(ctx context.Context, conn gotann.Connection, sessionID uint) ([]*Ticket, error)
	CheckIn(ctx context.Context, conn gotann.Connection, id uint) error
}

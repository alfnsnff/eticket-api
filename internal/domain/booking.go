package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type Booking struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	OrderID         string    `gorm:"column:order_id;type:varchar(64);not null;uniqueIndex"` // Business order ID
	ReferenceNumber *string   `gorm:"column:reference_number;"`
	ScheduleID      uint      `gorm:"column:schedule_id;not null;index;"`
	IDType          string    `gorm:"column:id_type;type:varchar(24);not null"`
	IDNumber        string    `gorm:"column:id_number;type:varchar(24);not null"`
	CustomerName    string    `gorm:"column:customer_name;type:varchar(32);not null"`
	CustomerAge     int       `gorm:"column:customer_age;not null"`
	CustomerGender  string    `gorm:"column:customer_gender;type:varchar(24);not null"`
	PhoneNumber     string    `gorm:"column:phone_number;type:varchar(14);not null"`
	Email           string    `gorm:"column:email;not null"`
	Status          string    `gorm:"column:status;type:varchar(24);not null;index"`
	ExpiresAt       time.Time `gorm:"column:expires_at;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null"`

	Tickets  []Ticket `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID"`
}

func (b *Booking) TableName() string {
	return "booking"
}

type BookingRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Booking) error
	InsertBulk(ctx context.Context, conn gotann.Connection, bookings []*Booking) error
	Update(ctx context.Context, conn gotann.Connection, entity *Booking) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, bookings []*Booking) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Booking) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Booking, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Booking, error)
	FindByOrderID(ctx context.Context, conn gotann.Connection, id string) (*Booking, error)
}

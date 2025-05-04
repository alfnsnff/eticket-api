package entity

import "time"

// Ticket represents an individual ticket slot for a specific class on a schedule.
// It tracks its status through the booking process and links to a confirmed booking.
type Ticket struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Links to Schedule and Class (defines what the ticket is for)
	ScheduleID uint `gorm:"not null;index;" json:"schedule_id"` // Foreign key to the Schedule (trip)
	ClassID    uint `gorm:"not null;index;" json:"class_id"`    // Foreign key to the Class

	// Link to the Booking transaction (NULLABLE until confirmed)
	BookingID *uint `gorm:"index" json:"booking_id"` // Foreign key to the Booking entity (NULL when pending)

	// Status and Timeout Management
	Status    string    `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'pending_data_entry', 'pending_payment', 'confirmed', 'cancelled'
	ClaimedAt time.Time `gorm:"not null" json:"claimed_at"`              // Timestamp when the ticket status became 'pending_data_entry'
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`              // Timeout timestamp for pending status

	// Pricing
	Price float32 `gorm:"not null" json:"price"` // Stored price at the time of claim/booking

	// Passenger Data (Filled after claiming, NULLABLE initially)
	PassengerName *string `json:"passenger_name"` // Pointer to string to allow NULL
	IDType        *string `json:"id_type"`        // Pointer to string to allow NULL (e.g., Passport, ID Card)
	IDNumber      *string `json:"ID_number"`      // Pointer to string to allow NULL
	// Add other passenger-specific fields here (e.g., DateOfBirth *time.Time, Gender *string)

	// Timestamps
	DataFilledAt     *time.Time `json:"data_filled_at"`     // Timestamp when passenger data was submitted (Pointer to time.Time to allow NULL)
	BookingTimestamp *time.Time `json:"booking_timestamp"`  // Timestamp when status became 'confirmed' (Pointer to time.Time to allow NULL)
	CreatedAt        time.Time  `json:"created_created_at"` // Standard creation timestamp
	UpdatedAt        time.Time  `json:"updated_at"`         // Standard update timestamp

	// Optional: Seat Number (if applicable and assigned)
	SeatNumber *string `json:"seat_number"` // Pointer to string to allow NULL

	// Relations (GORM tags for associations - optional in entity struct but helpful)
	// Booking  *Booking  `gorm:"foreignKey:BookingID"`  // Belongs to one Booking (when BookingID is not NULL)
	// Schedule *Schedule `gorm:"foreignKey:ScheduleID"` // Belongs to one Schedule
	// Class    *Class    `gorm:"foreignKey:ClassID"`    // Belongs to one Class
}

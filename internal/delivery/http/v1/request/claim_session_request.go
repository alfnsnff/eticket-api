package request

import (
	"time"
)

type ClaimSessionScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type ClaimSessionScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type ClaimSessionSchedule struct {
	ID                uint                       `json:"id"`
	Ship              ClaimSessionScheduleShip   `json:"ship"`
	DepartureHarbor   ClaimSessionScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor     ClaimSessionScheduleHarbor `json:"arrival_harbor"`
	DepartureDatetime time.Time                  `json:"departure_datetime"`
	ArrivalDatetime   time.Time                  `json:"arrival_datetime"`
}

type ClaimSessionItemClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type ClaimSessionItem struct {
	ClassID  uint                  `json:"class_id"`
	Class    ClaimSessionItemClass `json:"class"`
	Subtotal float64               `json:"subtotal"` // Total price for this class (Quantity * Price)
	Quantity int                   `json:"quantity"`
}

type TESTLockClaimSessionRequest struct {
	ScheduleID uint               `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionItem `json:"items"`       // List of classes and quantities requested
}

type TESTClaimSessionTicketDataEntry struct {
	ClassID         uint    `json:"class_id"` // The class ID for the ticket
	PassengerName   string  `json:"passenger_name"`
	IDType          string  `json:"id_type"`
	IDNumber        string  `json:"id_number"`
	PassengerAge    int     `json:"passenger_age"`
	PassengerGender string  `json:"passenger_gender"`
	Address         string  `json:"address"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
}

type TESTDataEntryClaimSessionRequest struct {
	CustomerName   string                            `json:"customer_name"`
	CustomerAge    int                               `json:"customer_age"`
	CustomerGender string                            `json:"customer_gender"`
	IDType         string                            `json:"id_type"`
	IDNumber       string                            `json:"id_number"`
	PhoneNumber    string                            `json:"phone_number"`
	Email          string                            `json:"email"`
	BirthDate      time.Time                         `json:"birth_date"`
	TicketData     []TESTClaimSessionTicketDataEntry `json:"ticket_data"`
	PaymentMethod  string                            `json:"payment_method"`
}

type ClaimSessionResponse struct {
	ID          uint                 `json:"id"`
	SessionID   string               `json:"session_id"`
	Schedule    ClaimSessionSchedule `json:"schedule"`
	Status      string               `json:"status"`
	ClaimItems  []ClaimSessionItem   `json:"claim_items"`
	TotalAmount float64              `json:"total_amount"`
	ExpiresAt   time.Time            `json:"expires_at"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

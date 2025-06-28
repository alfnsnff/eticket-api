package request

import "time"

type TESTLockClaimSessionRequest struct {
	ScheduleID uint               `json:"schedule_id"` // The schedule the user wants tickets for
	Items      []ClaimSessionItem `json:"items"`       // List of classes and quantities requested
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

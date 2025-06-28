package request

import "time"

type CreateClaimItemRequest struct {
	ClaimSessionID uint `json:"claim_session_id" validate:"required"`
	ClassID        uint `json:"class_id" validate:"required"`
	Quantity       int  `json:"quantity" validate:"required,min=1"`
}

type UpdateClaimItemRequest struct {
	ID             uint `json:"id" validate:"required"`
	ClaimSessionID uint `json:"claim_session_id" validate:"required"`
	ClassID        uint `json:"class_id" validate:"required"`
	Quantity       int  `json:"quantity" validate:"required,min=1"`
}

type ClaimItemResponse struct {
	ID             uint      `json:"id"`
	ClaimSessionID uint      `json:"claim_session_id"`
	ClassID        uint      `json:"class_id"`
	Quantity       int       `json:"quantity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

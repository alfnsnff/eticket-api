package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Class domain to ReadClassResponse model
func ClaimItemToResponse(claimItem *domain.ClaimItem) *model.ReadClaimItemResponse {
	return &model.ReadClaimItemResponse{
		ID:             claimItem.ID,
		ClaimSessionID: claimItem.ClaimSessionID,
		ClassID:        claimItem.ClassID,
		Quantity:       claimItem.Quantity,
		CreatedAt:      claimItem.CreatedAt,
		UpdatedAt:      claimItem.UpdatedAt,
	}
}

func ClaimItemFromRequest(req *model.WriteClaimItemRequest) *domain.ClaimItem {
	return &domain.ClaimItem{
		ClaimSessionID: req.ClaimSessionID,
		ClassID:        req.ClassID,
		Quantity:       req.Quantity,
	}
}

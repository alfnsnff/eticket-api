package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map ClaimSession domain to ReadClaimSessionResponse model
func TESTClaimSessionToResponse(session *domain.ClaimSession) *model.TESTReadClaimSessionResponse {
	claimItems := make([]model.ClaimSessionItem, len(session.ClaimItems))
	for i, item := range session.ClaimItems {
		claimItems[i] = model.ClaimSessionItem{
			ClassID: item.Class.ID,
			Class: model.ClaimSessionItemClass{
				ID:        item.Class.ID,
				ClassName: item.Class.ClassName,
				Type:      item.Class.Type,
			},
			Subtotal: item.Subtotal, // Assuming Subtotal is set in ClaimItem
			Quantity: item.Quantity,
		}
	}

	return &model.TESTReadClaimSessionResponse{
		ID:        session.ID,
		SessionID: session.SessionID,
		Status:    session.Status,
		Schedule: model.ClaimSessionSchedule{
			ID: session.Schedule.ID,
			Ship: model.ClaimSessionScheduleShip{
				ID:       session.Schedule.Ship.ID,
				ShipName: session.Schedule.Ship.ShipName,
			},
			DepartureHarbor: model.ClaimSessionScheduleHarbor{
				ID:         session.Schedule.DepartureHarbor.ID,
				HarborName: session.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.ClaimSessionScheduleHarbor{
				ID:         session.Schedule.ArrivalHarbor.ID,
				HarborName: session.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: session.Schedule.DepartureDatetime,
			ArrivalDatetime:   session.Schedule.ArrivalDatetime,
		},
		ExpiresAt:  session.ExpiresAt,
		ClaimItems: claimItems,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
	}
}

func ClaimSessionFromRequest(req *model.TESTWriteClaimSessionRequest) *domain.ClaimSession {
	claimItems := make([]domain.ClaimItem, len(req.Items))
	for i, item := range req.Items {
		claimItems[i] = domain.ClaimItem{
			ClassID:  item.ClassID,
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
		}
	}
	return &domain.ClaimSession{
		ScheduleID: req.ScheduleID,
		ClaimItems: claimItems,
	}
}

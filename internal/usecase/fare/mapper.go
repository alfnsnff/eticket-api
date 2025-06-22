package fare

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Fare domain to ReadFareResponse model
func FareToResponse(fare *domain.Fare) *model.ReadFareResponse {
	return &model.ReadFareResponse{
		ID: fare.ID,
		Route: model.FareRoute{
			ID: fare.Route.ID,
			DepartureHarbor: model.FareRouteHarbor{
				ID:         fare.Route.DepartureHarbor.ID,
				HarborName: fare.Route.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.FareRouteHarbor{
				ID:         fare.Route.ArrivalHarbor.ID,
				HarborName: fare.Route.ArrivalHarbor.HarborName,
			},
		},
		Manifest: model.FareManifest{
			ID: fare.Manifest.ID,
			Class: model.FareManifestClass{
				ID:        fare.Manifest.Class.ID,
				ClassName: fare.Manifest.Class.ClassName,
			},
			Ship: model.FareManifestShip{
				ID:       fare.Manifest.Ship.ID,
				ShipName: fare.Manifest.Ship.ShipName,
			},
			Capacity: fare.Manifest.Capacity,
		},
		TicketPrice: fare.TicketPrice,
		CreatedAt:   fare.CreatedAt,
		UpdatedAt:   fare.UpdatedAt,
	}
}

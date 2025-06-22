package fare

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Fare entity to ReadFareResponse model
func ToReadFareResponse(fare *entity.Fare) *model.ReadFareResponse {
	if fare == nil {
		return nil
	}

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

// Map a slice of Fare entities to ReadFareResponse models
func ToReadFareResponses(fares []*entity.Fare) []*model.ReadFareResponse {
	responses := make([]*model.ReadFareResponse, len(fares))
	for i, fare := range fares {
		responses[i] = ToReadFareResponse(fare)
	}
	return responses
}

// Map WriteFareRequest model to Fare entity
func FromWriteFareRequest(request *model.WriteFareRequest) *entity.Fare {
	return &entity.Fare{
		RouteID:     request.RouteID,
		ManifestID:  request.ManifestID,
		TicketPrice: request.TicketPrice,
	}
}

// Map UpdateFareRequest model to Fare entity
func FromUpdateFareRequest(request *model.UpdateFareRequest, fare *entity.Fare) {
	fare.RouteID = request.RouteID
	fare.ManifestID = request.ManifestID
	fare.TicketPrice = request.TicketPrice
}

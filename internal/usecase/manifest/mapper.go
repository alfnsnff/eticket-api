package manifest

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Manifest domain to ReadManifestResponse model
func ManifestToResponse(manifest *domain.Manifest) *model.ReadManifestResponse {
	return &model.ReadManifestResponse{
		ID: manifest.ID,
		Class: model.ManifestClass{
			ID:        manifest.Class.ID,
			ClassName: manifest.Class.ClassName,
			Type:      manifest.Class.Type,
		},
		Ship: model.ManifestShip{
			ID:       manifest.Ship.ID,
			ShipName: manifest.Ship.ShipName,
			ShipType: manifest.Ship.ShipType,
		},
		Capacity:  manifest.Capacity,
		CreatedAt: manifest.CreatedAt,
		UpdatedAt: manifest.UpdatedAt,
	}
}

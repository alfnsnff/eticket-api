package manifest

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Manifest entity to ReadManifestResponse model
func ToReadManifestResponse(manifest *entity.Manifest) *model.ReadManifestResponse {
	if manifest == nil {
		return nil
	}

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

// Map a slice of Manifest entities to ReadManifestResponse models
func ToReadManifestResponses(manifests []*entity.Manifest) []*model.ReadManifestResponse {
	responses := make([]*model.ReadManifestResponse, len(manifests))
	for i, manifest := range manifests {
		responses[i] = ToReadManifestResponse(manifest)
	}
	return responses
}

// Map WriteManifestRequest model to Manifest entity
func FromWriteManifestRequest(request *model.WriteManifestRequest) *entity.Manifest {
	return &entity.Manifest{
		ShipID:   request.ShipID,
		ClassID:  request.ClassID,
		Capacity: request.Capacity,
	}
}

// Map UpdateManifestRequest model to Manifest entity
func FromUpdateManifestRequest(request *model.UpdateManifestRequest, manifest *entity.Manifest) {
	manifest.ShipID = request.ShipID
	manifest.ClassID = request.ClassID
	manifest.Capacity = request.Capacity
}

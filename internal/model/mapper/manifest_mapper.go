package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToManifestModel(manifest *entities.Manifest) *model.ReadManifestResponse {
	response := new(model.ReadManifestResponse)
	copier.Copy(&response, &manifest) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToManifestsModel(manifests []*entities.Manifest) []*model.ReadManifestResponse {
	responses := []*model.ReadManifestResponse{}
	for _, manifest := range manifests {
		responses = append(responses, ToManifestModel(manifest))
	}
	return responses
}

func ToManifestEntity(request *model.WriteManifestRequest) *entities.Manifest {
	manifest := new(entities.Manifest)
	copier.Copy(&manifest, &request)
	return manifest
}

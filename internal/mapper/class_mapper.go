package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Class domain to ReadClassResponse model
func ClassToResponse(class *domain.Class) *model.ReadClassResponse {
	return &model.ReadClassResponse{
		ID:         class.ID,
		ClassName:  class.ClassName,
		ClassAlias: class.ClassAlias,
		Type:       class.Type,
		CreatedAt:  class.CreatedAt,
		UpdatedAt:  class.UpdatedAt,
	}
}

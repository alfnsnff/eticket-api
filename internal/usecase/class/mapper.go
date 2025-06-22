package class

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Class entity to ReadClassResponse model
func ToReadClassResponse(class *entity.Class) *model.ReadClassResponse {
	if class == nil {
		return nil
	}

	return &model.ReadClassResponse{
		ID:         class.ID,
		ClassName:  class.ClassName,
		ClassAlias: class.ClassAlias,
		Type:       class.Type,
		CreatedAt:  class.CreatedAt,
		UpdatedAt:  class.UpdatedAt,
	}
}

// Map a slice of Class entities to ReadClassResponse models
func ToReadClassResponses(classes []*entity.Class) []*model.ReadClassResponse {
	responses := make([]*model.ReadClassResponse, len(classes))
	for i, class := range classes {
		responses[i] = ToReadClassResponse(class)
	}
	return responses
}

// Map WriteClassRequest model to Class entity
func FromWriteClassRequest(request *model.WriteClassRequest) *entity.Class {
	return &entity.Class{
		ClassName:  request.ClassName,
		Type:       request.Type,
		ClassAlias: request.ClassAlias,
	}
}

// Map UpdateClassRequest model to Class entity
func FromUpdateClassRequest(request *model.UpdateClassRequest, class *entity.Class) {
	class.ClassName = request.ClassName
	class.Type = request.Type
	class.ClassAlias = request.ClassAlias
}

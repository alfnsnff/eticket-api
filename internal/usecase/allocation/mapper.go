package allocation

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Allocation entity to ReadAllocationResponse model
func ToReadAllocationResponse(allocation *entity.Allocation) *model.ReadAllocationResponse {
	if allocation == nil {
		return nil
	}

	return &model.ReadAllocationResponse{
		ID:         allocation.ID,
		ScheduleID: allocation.ScheduleID,
		Class: model.AllocationClass{
			ID:        allocation.Class.ID,
			ClassName: allocation.Class.ClassName,
			Type:      allocation.Class.Type,
		},
		Quota:     allocation.Quota,
		CreatedAt: allocation.CreatedAt,
		UpdatedAt: allocation.UpdatedAt,
	}
}

// Map a slice of Allocation entities to ReadAllocationResponse models
func ToReadAllocationResponses(allocations []*entity.Allocation) []*model.ReadAllocationResponse {
	responses := make([]*model.ReadAllocationResponse, len(allocations))
	for i, allocation := range allocations {
		responses[i] = ToReadAllocationResponse(allocation)
	}
	return responses
}

// Map WriteAllocationRequest model to Allocation entity
func FromWriteAllocationRequest(request *model.WriteAllocationRequest) *entity.Allocation {
	return &entity.Allocation{
		ScheduleID: request.ScheduleID,
		ClassID:    request.ClassID,
		Quota:      request.Quota,
	}
}

// Map UpdateAllocationRequest model to Allocation entity
func FromUpdateAllocationRequest(request *model.UpdateAllocationRequest, allocation *entity.Allocation) {
	allocation.ScheduleID = request.ScheduleID
	allocation.ClassID = request.ClassID
	allocation.Quota = request.Quota
}

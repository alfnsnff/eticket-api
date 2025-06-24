package allocation

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

func AllocationToResponse(allocation *domain.Allocation) *model.ReadAllocationResponse {
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

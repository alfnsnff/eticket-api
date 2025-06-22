package model

import "time"

type AllocationClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

// AllocationDTO represents a Allocation.
type ReadAllocationResponse struct {
	ID         uint            `json:"id"`
	ScheduleID uint            `json:"schedule_id"` // Foreign key
	Class      AllocationClass `json:"class"`       // Foreign key
	Quota      int             `json:"quota"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteAllocationRequest struct {
	ScheduleID uint `json:"schedule_id" validate:"required,gt=0"` // must be > 0
	ClassID    uint `json:"class_id" validate:"required,gt=0"`    // must be > 0
	Quota      int  `json:"quota" validate:"required,gte=0"`      // must be ≥ 0
}

// AllocationDTO represents a Allocation.
type UpdateAllocationRequest struct {
	ID         uint `json:"id" validate:"required,gt=0"`          // must be > 0
	ScheduleID uint `json:"schedule_id" validate:"required,gt=0"` // must be > 0
	ClassID    uint `json:"class_id" validate:"required,gt=0"`    // must be > 0
	Quota      int  `json:"quota" validate:"required,gte=0"`      // must be ≥ 0
}

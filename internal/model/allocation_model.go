package model

import "time"

// AllocationDTO represents a Allocation.
type ReadAllocationResponse struct {
	ID         uint      `json:"id"`
	ScheduleID uint      `json:"schedule_id"` // Foreign key
	ClassID    uint      `json:"class_id"`    // Foreign key
	Quota      int       `json:"quota"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteAllocationRequest struct {
	ScheduleID uint `json:"schedule_id"` // Foreign key
	ClassID    uint `json:"class_id"`    // Foreign key
	Quota      int  `json:"quota"`
}

// AllocationDTO represents a Allocation.
type UpdateAllocationRequest struct {
	ID         uint `json:"id"`
	ScheduleID uint `json:"schedule_id"` // Foreign key
	ClassID    uint `json:"class_id"`    // Foreign key
	Quota      int  `json:"quota"`
}

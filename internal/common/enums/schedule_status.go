package enum

// ScheduleStatus represents the status of a claim session
type ScheduleStatus int

const (
	ScheduleActive ScheduleStatus = iota
	ScheduleFinished
	ScheduleCancelled
)

func (ss ScheduleStatus) String() string {
	switch ss {
	case ScheduleActive:
		return "SCHEDULED"
	case ScheduleFinished:
		return "FINISHED"
	case ScheduleCancelled:
		return "CANCELLED"

	default:
		return "UNKNOWN"
	}
}

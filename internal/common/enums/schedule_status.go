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

// Get statuses that are confirmed and always occupy seats (ignore expiration)
func GetFinisihedScheduleStatuses() []string {
	return []string{
		ClaimSessionSuccess.String(),
	}
}

func GetActiveScheduleStatuses() []string {
	return []string{
		ClaimSessionSuccess.String(),
	}
}

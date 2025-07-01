package enum

// ScheduleStatus represents the status of a claim session
type HarborStatus int

const (
	HarborActive HarborStatus = iota
	HarborInactive
)

func (hs HarborStatus) String() string {
	switch hs {
	case HarborActive:
		return "ACTIVE"
	case HarborInactive:
		return "INACTIVE"

	default:
		return "UNKNOWN"
	}
}

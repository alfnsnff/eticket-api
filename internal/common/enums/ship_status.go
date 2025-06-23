package enum

// ShipStatus represents the status of a claim session
type ShipStatus int

const (
	ShipActive ShipStatus = iota
	ShipOnMaintenance
	ShipInactive
)

func (ss ShipStatus) String() string {
	switch ss {
	case ShipActive:
		return "ACTIVE"
	case ShipOnMaintenance:
		return "DOCKED"
	case ShipInactive:
		return "INACTIVE"

	default:
		return "unknown"
	}
}

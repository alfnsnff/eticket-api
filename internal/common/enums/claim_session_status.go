package enum

// ClaimSessionStatus represents the status of a claim session
type ClaimSessionStatus int

const (
	ClaimSessionFailed ClaimSessionStatus = iota
	ClaimSessionPending
	ClaimSessionSuccess
	ClaimSessionCancelled
)

func (css ClaimSessionStatus) String() string {
	switch css {
	case ClaimSessionFailed:
		return "FAILED"
	case ClaimSessionPending:
		return "PENDING"
	case ClaimSessionSuccess:
		return "RESERVED"
	case ClaimSessionCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

// Get statuses that are confirmed and always occupy seats (ignore expiration)
func GetSuccessClaimSessionStatuses() []string {
	return []string{
		ClaimSessionSuccess.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetPendingClaimSessionStatuses() []string {
	return []string{
		ClaimSessionPending.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetFailedClaimSessionStatuses() []string {
	return []string{
		ClaimSessionFailed.String(),
	}
}

package enum

// ClaimSessionStatus represents the status of a claim session
type ClaimSessionStatus int

const (
	ClaimSessionFailed ClaimSessionStatus = iota
	ClaimSessionPendingData
	ClaimSessionPendingPayment
	ClaimSessionPendingTransaction
	ClaimSessionSuccess
	ClaimSessionCancelled
)

func (css ClaimSessionStatus) String() string {
	switch css {
	case ClaimSessionFailed:
		return "failed"
	case ClaimSessionPendingData:
		return "data"
	case ClaimSessionPendingPayment:
		return "payment"
	case ClaimSessionPendingTransaction:
		return "transaction"
	case ClaimSessionSuccess:
		return "success"
	case ClaimSessionCancelled:
		return "cancelled"
	default:
		return "unknown"
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
		ClaimSessionPendingData.String(),
		ClaimSessionPendingPayment.String(),
		ClaimSessionPendingTransaction.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetFailedClaimSessionStatuses() []string {
	return []string{
		ClaimSessionFailed.String(),
	}
}

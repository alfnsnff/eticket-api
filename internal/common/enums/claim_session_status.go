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
		return "FAILED"
	case ClaimSessionPendingData:
		return "DATA_PENDING"
	case ClaimSessionPendingPayment:
		return "PAYMENT_PENDING"
	case ClaimSessionPendingTransaction:
		return "TRANSACTION_PENDING"
	case ClaimSessionSuccess:
		return "SCUCCESS"
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

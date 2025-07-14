package enum

// BookingStatus represents the status of a claim session
type BookingStatus int

const (
	BookingPaid BookingStatus = iota
	BookingUnpaid
	BookingExpired
	BookingRefund
)

func (css BookingStatus) String() string {
	switch css {
	case BookingPaid:
		return "PAID"
	case BookingUnpaid:
		return "UNPAID"
	case BookingExpired:
		return "EXPIRED"
	case BookingRefund:
		return "REFUND"
	default:
		return "FAILED"
	}
}

// Get statuses that are confirmed and always occupy seats (ignore expiration)
func GetSuccessBookingStatuses() []string {
	return []string{
		BookingPaid.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetPendingBookingStatuses() []string {
	return []string{
		BookingUnpaid.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetFailedBookingStatuses() []string {
	return []string{
		BookingExpired.String(),
	}
}

// Get statuses that are pending and respect expiration
func GetInvalidBookingStatuses() []string {
	return []string{
		BookingExpired.String(),
		BookingUnpaid.String(),
		BookingRefund.String(),
	}
}

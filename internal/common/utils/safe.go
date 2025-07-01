package utils

import "time"

func SafeStringDeref(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func SafeUintDeref(ptr *uint) uint {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func SafeIntDeref(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func SafeFloat64Deref(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0.0
}

func SafeFloat32Deref(ptr *float32) float32 {
	if ptr != nil {
		return *ptr
	}
	return 0.0
}

func SafeTimeDeref(ptr *time.Time) time.Time {
	if ptr != nil {
		return *ptr
	}
	return time.Time{}
}

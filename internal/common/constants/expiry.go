package constant

import "time"

const (
	AccessTokenExpiry  = 7 * 24 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour // 7 days
)

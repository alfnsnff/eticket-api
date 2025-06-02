package code

import (
	"crypto/rand"
	"fmt"
	"time"
)

func randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

func GenerateOrderID(routeCode, shipCode string, date time.Time) string {
	dateStr := date.Format("20060102")
	return fmt.Sprintf("TH-%s-%s-%s-%s", randomString(5), routeCode, shipCode, dateStr)
}

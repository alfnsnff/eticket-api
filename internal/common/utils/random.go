package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

func random(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
func GenerateOrderID(routeCode, shipCode string, date time.Time) string {
	timestamp := date.Format("20060102150405") // Termasuk detik
	return fmt.Sprintf("ID-%s-%s%s-%s", routeCode, shipCode, timestamp, random(5))
}

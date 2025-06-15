package utils

import (
	crypt "crypto/rand"
	"fmt"
	math "math/rand"
	"time"
)

// GenerateOrderID creates a unique order ID with a shuffled timestamp.
func GenerateOrderID(routeCode, shipCode string, date time.Time) string {
	timestamp := date.Format("20060102150405") // e.g. "20250615173542"
	return fmt.Sprintf("ID-%s-%s%s-%s", routeCode, shipCode, ShuffleString(timestamp), Random(5))
}

func Random(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345678901234567890123456789"
	b := make([]byte, n)
	_, _ = crypt.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

// shuffleString randomly shuffles the characters in a string.
func ShuffleString(s string) string {
	runes := []rune(s)
	math.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

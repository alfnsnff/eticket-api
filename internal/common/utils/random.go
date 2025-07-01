package utils

import (
	crypt "crypto/rand"
	"fmt"
	math "math/rand"
	"time"
)

func GenerateOrderID(departure string) string {
	now := time.Now()
	day := now.Format("02")  // Tanggal (DD)
	hour := now.Format("15") // Jam (HH)

	// Get last 3 digits of millisecond timestamp
	msPart := now.UnixMilli() % 1000
	msStr := fmt.Sprintf("%03d", msPart)

	random := Random(4) // A-Z0-9, 4 chars

	return fmt.Sprintf("%s%s%s%s%s", departure, day, hour, msStr, random)
}

func Random(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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

func GenerateTicketReferenceID() string {
	now := time.Now().Format("060102") // YYMMDD
	random := Random(3)                // 6-char alphanumeric
	return fmt.Sprintf("T-%s-%s", now, random)
}

package utils

import (
	crypt "crypto/rand"
	"encoding/hex"
	"fmt"
	math "math/rand"
	"time"
)

// GenerateOrderID generates a highly unique order ID
func GenerateOrderID(departure string) string {
	now := time.Now().UTC()
	timePart := now.Format("20060102150405") // yyyyMMddHHmmss
	nanoPart := fmt.Sprintf("%09d", now.Nanosecond())

	uuidPart := generateRandomHex(8) // 8 bytes = 16 hex chars

	return fmt.Sprintf("%s-%s%s-%s", departure, timePart, nanoPart, uuidPart)
}

// GenerateTicketReferenceID creates a highly unique ticket reference ID
func GenerateTicketReferenceID() string {
	now := time.Now().UTC()
	datePart := now.Format("060102")                      // YYMMDD
	nanoPart := fmt.Sprintf("%09d", now.Nanosecond())[:6] // 6-digit nanosec
	randPart := generateRandomHex(4)                      // 4 bytes = 8 hex chars

	return fmt.Sprintf("T%s-%s%s", datePart, nanoPart, randPart)
}

// generateRandomHex returns a securely generated random hex string of n bytes
func generateRandomHex(n int) string {
	b := make([]byte, n)
	if _, err := crypt.Read(b); err != nil {
		panic("failed to generate secure random ID")
	}
	return hex.EncodeToString(b)
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

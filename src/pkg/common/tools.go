package common

import (
	"ledger/src/pkg/uuid"
	"time"
	"unicode"
)

// CreateTimestamp generates a current timestamp in milliseconds
// Returns the number of milliseconds since the Unix epoch (January 1, 1970 UTC)
func CreateTimestamp() int64 {
	// Get current time in nanoseconds and convert to milliseconds
	t := time.Now().UnixNano() / 1e6
	return t
}

// CreateUuid generates a new UUID v4
// Returns a string representation of a version 4 UUID
func CreateUuid() string {
	return uuid.V4UUID()
}

// IsValid32UUID checks if a string is a valid 32-character UUID (without hyphens)
// Returns true if the string is exactly 32 characters long
// and contains only valid hexadecimal characters (0-9, a-f, A-F)
func IsValid32UUID(uuid string) bool {
	// First check if length is exactly 32 characters
	if len(uuid) != 32 {
		return false
	}

	// Check each character is a valid hexadecimal digit
	for _, c := range uuid {
		if !unicode.Is(unicode.Hex_Digit, c) {
			return false
		}
	}

	return true
}

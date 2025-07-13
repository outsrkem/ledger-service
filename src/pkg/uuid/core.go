package uuid

import (
	"encoding/binary"
	"fmt"
)

// UUID represents a 128-bit UUID (Universally Unique Identifier)
type UUID [16]byte

// String returns the UUID in standard 8-4-4-4-12 format
func (u UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(u[0:4]),
		binary.BigEndian.Uint16(u[4:6]),
		binary.BigEndian.Uint16(u[6:8]),
		binary.BigEndian.Uint16(u[8:10]),
		u[10:16])
}

// NoHyphenString returns the UUID without hyphens
func (u UUID) NoHyphenString() string {
	return fmt.Sprintf("%08x%04x%04x%04x%012x",
		binary.BigEndian.Uint32(u[0:4]),
		binary.BigEndian.Uint16(u[4:6]),
		binary.BigEndian.Uint16(u[6:8]),
		binary.BigEndian.Uint16(u[8:10]),
		u[10:16])
}

// Version returns the UUID version number (1-5)
func (u UUID) Version() uint8 {
	return (u[6] >> 4) & 0x0F
}

// Variant returns the UUID variant
func (u UUID) Variant() uint8 {
	return (u[8] >> 6) & 0x03
}

// IsValid checks if the UUID conforms to RFC 4122 standards
func (u UUID) IsValid() bool {
	version := u.Version()
	variant := u.Variant()

	// Valid versions are 1-5, valid RFC 4122 variant is 2
	return (version >= 1 && version <= 5) && variant == 2
}

package uuid

import (
	"crypto/rand"
	"fmt"
)

// Predefined namespaces as specified in RFC 4122
var (
	NamespaceDNS  = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceURL  = UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceOID  = UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceX500 = UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
)

// NewV4 generates a version 4 UUID (random-based)
func NewV4() (UUID, error) {
	var u UUID
	if _, err := rand.Read(u[:]); err != nil {
		return UUID{}, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Set version to 4 (0100 in the most significant 4 bits of the 7th byte)
	u[6] = (u[6] & 0x0F) | 0x40

	// Set variant to RFC 4122 (10 in the most significant 2 bits of the 9th byte)
	u[8] = (u[8] & 0x3F) | 0x80

	return u, nil
}

// NewV4Batch generates multiple version 4 UUIDs in batch
func NewV4Batch(count int) ([]UUID, error) {
	if count <= 0 {
		return nil, fmt.Errorf("invalid count: %d", count)
	}

	batch := make([]UUID, count)
	data := make([]byte, 16*count)

	if _, err := rand.Read(data); err != nil {
		return nil, fmt.Errorf("failed to generate batch random bytes: %w", err)
	}

	for i := 0; i < count; i++ {
		copy(batch[i][:], data[i*16:(i+1)*16])

		// Set version and variant for each UUID
		batch[i][6] = (batch[i][6] & 0x0F) | 0x40
		batch[i][8] = (batch[i][8] & 0x3F) | 0x80
	}

	return batch, nil
}

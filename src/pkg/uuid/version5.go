package uuid

import "crypto/sha1"

// NewV5 generates a version 5 UUID (SHA-1 based)
func NewV5(namespace UUID, name []byte) UUID {
	hash := sha1.New()
	hash.Write(namespace[:])
	hash.Write(name)
	digest := hash.Sum(nil)

	var uuid UUID
	copy(uuid[:], digest[:16])

	// Set version to 5 (0101 in the most significant 4 bits of the 7th byte)
	uuid[6] = (uuid[6] & 0x0F) | 0x50

	// Set variant to RFC 4122 (10 in the most significant 2 bits of the 9th byte)
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return uuid
}

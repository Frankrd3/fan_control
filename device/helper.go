package device

import "github.com/google/uuid"

// FormatGUIDAsString FormatAsString Format RFC 4122 UUID as user readable string (https://github.com/bougou/go-ipmi/blob/b7dd375594178ada365455ae8fcf34b8f4b290ca/cmd_get_device_guid.go#L40)
func FormatGUIDAsString(GUID [16]byte) string {
	// UUID Most Significant Bit (MSB)
	// https://datatracker.ietf.org/doc/html/draft-ietf-uuidrev-rfc4122bis-12
	rfc4122Msb := make([]byte, 16)

	for i := 0; i < 16; i++ {
		rfc4122Msb[i] = GUID[:][15-i]
	}

	u, err := uuid.FromBytes(rfc4122Msb)

	if err != nil {
		return ""
	}

	return u.String()
}

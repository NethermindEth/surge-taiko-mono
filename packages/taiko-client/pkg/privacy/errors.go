package privacy

import (
	"errors"
	"fmt"
)

// ErrTruncated is returned when the blob payload is shorter than the minimum length
// required by its scheme header (e.g. less than nonce+tag bytes for AES).
var ErrTruncated = errors.New("privacy: blob payload truncated")

// ErrKeyMissing is returned when a blob requests a scheme whose decryption key is not
// configured on this component.
var ErrKeyMissing = errors.New("privacy: key not configured for scheme")

// ErrInvalidKey is returned when key bytes have the wrong length or invalid format.
var ErrInvalidKey = errors.New("privacy: invalid key bytes")

// ErrInvalidEphemeralPubkey is returned when the ECIES `pk_eph` bytes cannot be parsed.
var ErrInvalidEphemeralPubkey = errors.New("privacy: invalid ephemeral pubkey")

// ErrAEADFailed is returned on AES-GCM authentication failure (bad key, nonce, or tag).
var ErrAEADFailed = errors.New("privacy: AEAD authentication failed")

// UnknownSchemeError indicates the scheme byte is not in the registry.
type UnknownSchemeError struct {
	Scheme uint8
}

func (e *UnknownSchemeError) Error() string {
	return fmt.Sprintf("privacy: unknown scheme 0x%02x", e.Scheme)
}

// Package privacy implements driver-side decryption of realtime proposal blobs.
//
// Each privacy blob payload (the inner buffer carried by an EIP-4844 sidecar after
// the existing [version (1B)][size (3B BE)] framing) starts with a 1-byte scheme id
// that selects the cipher used. See PRIVACY_STACK.md at the repo root for the full
// byte-layout specification.
//
// Schemes:
//   - 0x00 = plaintext: payload is the compressed manifest verbatim.
//   - 0x01 = AES-256-GCM: a single shared symmetric key, fresh per-blob nonce on the wire.
//   - 0x02 = ECIES (secp256k1 + AES-GCM): for forced-inclusion blobs, encrypted to the
//     system's static public key by an external submitter.
package privacy

const (
	// SchemePlain indicates the rest of the payload is the compressed manifest verbatim.
	SchemePlain uint8 = 0x00
	// SchemeAES256GCM indicates AES-256-GCM with a shared symmetric key.
	SchemeAES256GCM uint8 = 0x01
	// SchemeECIESSecp256k1 indicates ECIES = secp256k1 ECDH ⊕ HKDF-SHA256 ⊕ AES-256-GCM.
	SchemeECIESSecp256k1 uint8 = 0x02
)

// Keys bundles the optional decryption keys consumed by Dispatch. Only the key required
// by a blob's actual scheme must be set; the others may be nil.
type Keys struct {
	// Symmetric is the 32-byte shared AES-256-GCM key used by Catalyst's normal proposals
	// (scheme 0x01).
	Symmetric []byte
	// FIPrivate is the 32-byte secp256k1 scalar (system FI private key, scheme 0x02).
	FIPrivate []byte
}

// Dispatch reads the scheme byte and routes to the right decoder.
//
// `payload` is the blob's inner buffer AFTER the [version][size] framing has been
// stripped — i.e. the bytes whose first element is the scheme id. For scheme 0x00,
// the function returns the rest verbatim; otherwise it decrypts under the matching key.
func Dispatch(payload []byte, keys Keys) ([]byte, error) {
	if len(payload) == 0 {
		return nil, ErrTruncated
	}
	scheme := payload[0]
	rest := payload[1:]
	switch scheme {
	case SchemePlain:
		out := make([]byte, len(rest))
		copy(out, rest)
		return out, nil
	case SchemeAES256GCM:
		if len(keys.Symmetric) == 0 {
			return nil, ErrKeyMissing
		}
		return aesDecrypt(rest, keys.Symmetric)
	case SchemeECIESSecp256k1:
		if len(keys.FIPrivate) == 0 {
			return nil, ErrKeyMissing
		}
		return eciesDecrypt(rest, keys.FIPrivate)
	default:
		return nil, &UnknownSchemeError{Scheme: scheme}
	}
}

// EciesInfo is the HKDF-SHA256 info string used to derive (key, nonce) from the
// ECDH shared secret in scheme 0x02. Submitter and system MUST use this exact string.
var EciesInfo = []byte("surge-fi-v1")

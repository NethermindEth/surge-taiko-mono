package privacy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/hkdf"
)

const (
	// pkLen is the length of a compressed secp256k1 pubkey in bytes (header + 32-byte x-coord).
	pkLen = 33

	// minEciesInnerLen is the minimum length of a scheme-0x02 inner payload (pk_eph + tag).
	minEciesInnerLen = pkLen + aesTagLen
)

// eciesDecrypt decrypts a scheme-0x02 inner payload `[pk_eph(33) || ct || tag(16)]`
// using the system private key `sk_sys` (32-byte secp256k1 scalar). Returns the
// plaintext compressed manifest.
func eciesDecrypt(inner []byte, skSys []byte) ([]byte, error) {
	if len(inner) < minEciesInnerLen {
		return nil, ErrTruncated
	}
	if len(skSys) != aesKeyLen {
		return nil, ErrInvalidKey
	}

	pkEphBytes := inner[:pkLen]
	ctAndTag := inner[pkLen:]

	pkEph, err := crypto.DecompressPubkey(pkEphBytes)
	if err != nil {
		return nil, ErrInvalidEphemeralPubkey
	}

	priv, err := crypto.ToECDSA(skSys)
	if err != nil {
		return nil, ErrInvalidKey
	}

	// ECDH: shared = sk_sys * pk_eph, take the X coordinate as 32-byte big-endian.
	curve := crypto.S256()
	x, _ := curve.ScalarMult(pkEph.X, pkEph.Y, priv.D.Bytes())

	shared := make([]byte, 32)
	xBytes := x.Bytes()
	copy(shared[32-len(xBytes):], xBytes)

	// HKDF-SHA256(salt=∅, ikm=shared, info="surge-fi-v1", L=44).
	keyAndNonce, err := hkdfExpand(shared, EciesInfo, aesKeyLen+aesNonceLen)
	if err != nil {
		return nil, ErrAEADFailed
	}
	kEph := keyAndNonce[:aesKeyLen]
	nonce := keyAndNonce[aesKeyLen:]

	block, err := aes.NewCipher(kEph)
	if err != nil {
		return nil, ErrInvalidKey
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrInvalidKey
	}
	pt, err := gcm.Open(nil, nonce, ctAndTag, nil)
	if err != nil {
		return nil, ErrAEADFailed
	}
	return pt, nil
}

// hkdfExpand performs HKDF-SHA256 with empty salt over `ikm` and returns `outLen` bytes.
func hkdfExpand(ikm, info []byte, outLen int) ([]byte, error) {
	r := hkdf.New(sha256.New, ikm, nil, info)
	out := make([]byte, outLen)
	if _, err := r.Read(out); err != nil {
		return nil, err
	}
	return out, nil
}

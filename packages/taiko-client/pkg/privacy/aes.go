package privacy

import (
	"crypto/aes"
	"crypto/cipher"
)

const (
	aesKeyLen   = 32
	aesNonceLen = 12
	aesTagLen   = 16

	// minAesInnerLen is the minimum length of a scheme-0x01 inner payload (nonce + tag).
	minAesInnerLen = aesNonceLen + aesTagLen
)

// aesDecrypt decrypts a scheme-0x01 inner payload `[nonce(12) || ct || tag(16)]` and
// returns the plaintext compressed manifest.
func aesDecrypt(inner []byte, key []byte) ([]byte, error) {
	if len(inner) < minAesInnerLen {
		return nil, ErrTruncated
	}
	if len(key) != aesKeyLen {
		return nil, ErrInvalidKey
	}

	nonce := inner[:aesNonceLen]
	ctAndTag := inner[aesNonceLen:]

	block, err := aes.NewCipher(key)
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

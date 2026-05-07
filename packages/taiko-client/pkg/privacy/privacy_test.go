package privacy

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/hkdf"
)

func TestDispatchPlaintext(t *testing.T) {
	t.Parallel()
	payload := []byte{SchemePlain, 'h', 'i'}
	out, err := Dispatch(payload, Keys{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(out, []byte("hi")) {
		t.Fatalf("got %q, want %q", out, "hi")
	}
}

func TestDispatchUnknownScheme(t *testing.T) {
	t.Parallel()
	payload := []byte{0xFF, 0xAA}
	_, err := Dispatch(payload, Keys{})
	var unknown *UnknownSchemeError
	if !errors.As(err, &unknown) || unknown.Scheme != 0xFF {
		t.Fatalf("got %v, want UnknownSchemeError(0xFF)", err)
	}
}

func TestDispatchTruncated(t *testing.T) {
	t.Parallel()
	_, err := Dispatch(nil, Keys{})
	if !errors.Is(err, ErrTruncated) {
		t.Fatalf("got %v, want ErrTruncated", err)
	}
}

func TestDispatchAesMissingKey(t *testing.T) {
	t.Parallel()
	_, err := Dispatch([]byte{SchemeAES256GCM, 0x00}, Keys{})
	if !errors.Is(err, ErrKeyMissing) {
		t.Fatalf("got %v, want ErrKeyMissing", err)
	}
}

func TestAesRoundtrip(t *testing.T) {
	t.Parallel()
	key := bytes.Repeat([]byte{0x42}, 32)
	nonce := bytes.Repeat([]byte{0x37}, 12)
	plaintext := []byte("compressed manifest goes here")

	// Encrypt the same way Catalyst would: nonce || ct || tag.
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	ctAndTag := gcm.Seal(nil, nonce, plaintext, nil)

	inner := append([]byte{}, nonce...)
	inner = append(inner, ctAndTag...)
	payload := append([]byte{SchemeAES256GCM}, inner...)

	out, err := Dispatch(payload, Keys{Symmetric: key})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(out, plaintext) {
		t.Fatalf("got %q, want %q", out, plaintext)
	}
}

func TestAesTamperedFails(t *testing.T) {
	t.Parallel()
	key := bytes.Repeat([]byte{0x42}, 32)
	nonce := bytes.Repeat([]byte{0x37}, 12)
	plaintext := []byte("abcdef")

	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	ctAndTag := gcm.Seal(nil, nonce, plaintext, nil)
	inner := append([]byte{}, nonce...)
	inner = append(inner, ctAndTag...)

	// Flip one ciphertext byte.
	inner[len(nonce)+1] ^= 0x01

	payload := append([]byte{SchemeAES256GCM}, inner...)
	_, err := Dispatch(payload, Keys{Symmetric: key})
	if !errors.Is(err, ErrAEADFailed) {
		t.Fatalf("got %v, want ErrAEADFailed", err)
	}
}

func TestEciesRoundtrip(t *testing.T) {
	t.Parallel()

	// System keypair.
	skSys, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("genkey: %v", err)
	}
	pkSysCompressed := crypto.CompressPubkey(&skSys.PublicKey)
	skSysBytes := crypto.FromECDSA(skSys)

	// Submitter side: ephemeral keypair + ECDH + HKDF + AES-GCM.
	skEph, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("genkey eph: %v", err)
	}
	pkEphCompressed := crypto.CompressPubkey(&skEph.PublicKey)

	pkSys, err := crypto.DecompressPubkey(pkSysCompressed)
	if err != nil {
		t.Fatalf("decompress pkSys: %v", err)
	}
	x, _ := crypto.S256().ScalarMult(pkSys.X, pkSys.Y, skEph.D.Bytes())
	shared := make([]byte, 32)
	xb := x.Bytes()
	copy(shared[32-len(xb):], xb)

	r := hkdf.New(sha256.New, shared, nil, EciesInfo)
	out := make([]byte, 44)
	if _, err := r.Read(out); err != nil {
		t.Fatalf("hkdf: %v", err)
	}
	kEph := out[:32]
	nonce := out[32:]

	plaintext := []byte("forced inclusion plaintext payload")
	block, _ := aes.NewCipher(kEph)
	gcm, _ := cipher.NewGCM(block)
	ctAndTag := gcm.Seal(nil, nonce, plaintext, nil)

	inner := append([]byte{}, pkEphCompressed...)
	inner = append(inner, ctAndTag...)
	payload := append([]byte{SchemeECIESSecp256k1}, inner...)

	got, err := Dispatch(payload, Keys{FIPrivate: skSysBytes})
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("got %q, want %q", got, plaintext)
	}
}

func TestEciesWrongSystemKeyFails(t *testing.T) {
	t.Parallel()

	skSys, _ := crypto.GenerateKey()
	pkSysCompressed := crypto.CompressPubkey(&skSys.PublicKey)
	skEph, _ := crypto.GenerateKey()
	pkEphCompressed := crypto.CompressPubkey(&skEph.PublicKey)

	pkSys, _ := crypto.DecompressPubkey(pkSysCompressed)
	x, _ := crypto.S256().ScalarMult(pkSys.X, pkSys.Y, skEph.D.Bytes())
	shared := make([]byte, 32)
	xb := x.Bytes()
	copy(shared[32-len(xb):], xb)

	r := hkdf.New(sha256.New, shared, nil, EciesInfo)
	out := make([]byte, 44)
	_, _ = r.Read(out)

	block, _ := aes.NewCipher(out[:32])
	gcm, _ := cipher.NewGCM(block)
	ctAndTag := gcm.Seal(nil, out[32:], []byte("x"), nil)

	inner := append([]byte{}, pkEphCompressed...)
	inner = append(inner, ctAndTag...)
	payload := append([]byte{SchemeECIESSecp256k1}, inner...)

	// Try to decrypt with a DIFFERENT system key — must fail.
	otherSk, _ := crypto.GenerateKey()
	_, err := Dispatch(payload, Keys{FIPrivate: crypto.FromECDSA(otherSk)})
	if !errors.Is(err, ErrAEADFailed) {
		t.Fatalf("got %v, want ErrAEADFailed", err)
	}
}

func TestAesRandomNoncesProducedByEncryptor(t *testing.T) {
	t.Parallel()
	// Smoke-test that this package's Dispatch works on an inner payload built with a
	// random nonce — what Catalyst will emit.
	key := bytes.Repeat([]byte{0x55}, 32)
	plaintext := []byte("hello")

	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("rand: %v", err)
	}
	ctAndTag := gcm.Seal(nil, nonce, plaintext, nil)
	inner := append([]byte{}, nonce...)
	inner = append(inner, ctAndTag...)
	payload := append([]byte{SchemeAES256GCM}, inner...)

	got, err := Dispatch(payload, Keys{Symmetric: key})
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("got %q, want %q", got, plaintext)
	}
}

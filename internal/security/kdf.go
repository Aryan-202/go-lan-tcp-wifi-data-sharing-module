package security

import (
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// DefaultAESKeySize is 32 bytes (256 bits) for AES-256
const DefaultAESKeySize = 32

// DeriveKey expands a raw shared secret into a uniform 32-byte AES key using HKDF-SHA256
func DeriveKey(sharedSecret []byte, salt []byte, info []byte) ([]byte, error) {
	if len(sharedSecret) == 0 {
		return nil, fmt.Errorf("shared secret cannot be empty")
	}

	if len(salt) == 0 {
		salt = []byte("GoShare-HKDF-Salt-v1")
	}

	if len(info) == 0 {
		info = []byte("GoShare-Session-Key")
	}

	hkdfReader := hkdf.New(sha256.New, sharedSecret, salt, info)
	derivedKey := make([]byte, DefaultAESKeySize)

	_, err := io.ReadFull(hkdfReader, derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key using HKDF: %w", err)
	}

	return derivedKey, nil
}

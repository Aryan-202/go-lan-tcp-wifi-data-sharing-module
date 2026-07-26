package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidKeySize is returned when key length is not 32 bytes for AES-256
var ErrInvalidKeySize = errors.New("AES-256 key must be exactly 32 bytes")

// ErrCiphertextTooShort is returned when ciphertext is shorter than GCM nonce size
var ErrCiphertextTooShort = errors.New("ciphertext is too short")

// ErrAuthenticationFailed is returned when AES-GCM tag validation fails (data tampered)
var ErrAuthenticationFailed = errors.New("decryption failed: data corrupted or tampered")

// Encrypt encrypts plaintext using AES-256-GCM with optional Additional Authenticated Data (AAD)
func Encrypt(key []byte, plaintext []byte, aad []byte) ([]byte, error) {
	if len(key) != DefaultAESKeySize {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %w", err)
	}

	// gcm.Seal(dst, nonce, plaintext, additionalData)
	// Passing nonce as dst prepends the 12-byte nonce to the encrypted output
	ciphertext := gcm.Seal(nonce, nonce, plaintext, aad)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM and verifies authentication tag
func Decrypt(key []byte, ciphertext []byte, aad []byte) ([]byte, error) {
	if len(key) != DefaultAESKeySize {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, aad)
	if err != nil {
		return nil, ErrAuthenticationFailed
	}

	return plaintext, nil
}

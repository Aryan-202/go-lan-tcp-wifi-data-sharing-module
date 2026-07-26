package security_test

import (
	"bytes"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

func TestEncryptDecrypt_Success(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := []byte("GoShare Secure Data Transfer Module Test Payload")
	aad := []byte("header-metadata")

	ciphertext, err := security.Encrypt(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := security.Decrypt(key, ciphertext, aad)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text does not match plaintext")
	}
}

func TestEncryptDecrypt_InvalidKeySize(t *testing.T) {
	invalidKey := []byte("short-key")
	plaintext := []byte("data")

	_, err := security.Encrypt(invalidKey, plaintext, nil)
	if err != security.ErrInvalidKeySize {
		t.Errorf("Expected ErrInvalidKeySize on Encrypt, got: %v", err)
	}

	_, err = security.Decrypt(invalidKey, plaintext, nil)
	if err != security.ErrInvalidKeySize {
		t.Errorf("Expected ErrInvalidKeySize on Decrypt, got: %v", err)
	}
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	key := make([]byte, 32)
	shortCiphertext := []byte("too-short")

	_, err := security.Decrypt(key, shortCiphertext, nil)
	if err != security.ErrCiphertextTooShort {
		t.Errorf("Expected ErrCiphertextTooShort, got: %v", err)
	}
}

func TestDecrypt_AuthenticationFailed(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("Authentic message")
	aad := []byte("header")

	ciphertext, err := security.Encrypt(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 1. Tamper ciphertext
	tamperedCiphertext := make([]byte, len(ciphertext))
	copy(tamperedCiphertext, ciphertext)
	tamperedCiphertext[len(tamperedCiphertext)-1] ^= 0xFF

	_, err = security.Decrypt(key, tamperedCiphertext, aad)
	if err != security.ErrAuthenticationFailed {
		t.Errorf("Expected ErrAuthenticationFailed for tampered ciphertext, got: %v", err)
	}

	// 2. Tamper AAD
	badAAD := []byte("tampered-header")
	_, err = security.Decrypt(key, ciphertext, badAAD)
	if err != security.ErrAuthenticationFailed {
		t.Errorf("Expected ErrAuthenticationFailed for bad AAD, got: %v", err)
	}
}

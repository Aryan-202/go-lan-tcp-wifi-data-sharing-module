package security_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

func TestX25519KeyExchangeAndHKDF(t *testing.T) {
	// 1. Generate KeyPairs for Initiator (Peer A) and Responder (Peer B)
	initiatorKeyPair, err := security.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Initiator GenerateKeyPair failed: %v", err)
	}

	responderKeyPair, err := security.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Responder GenerateKeyPair failed: %v", err)
	}

	// 2. Perform ECDH key exchange
	initiatorSecret, err := security.ComputeSharedSecret(initiatorKeyPair.PrivKey, responderKeyPair.PubKeyBytes())
	if err != nil {
		t.Fatalf("Initiator ComputeSharedSecret failed: %v", err)
	}

	responderSecret, err := security.ComputeSharedSecret(responderKeyPair.PrivKey, initiatorKeyPair.PubKeyBytes())
	if err != nil {
		t.Fatalf("Responder ComputeSharedSecret failed: %v", err)
	}

	if !bytes.Equal(initiatorSecret, responderSecret) {
		t.Fatalf("Computed raw secrets do not match!")
	}

	// 3. Derive AES-256 keys via HKDF-SHA256
	salt := []byte("test-salt-12345")
	info := []byte("goshare-session")

	initiatorKey, err := security.DeriveKey(initiatorSecret, salt, info)
	if err != nil {
		t.Fatalf("Initiator DeriveKey failed: %v", err)
	}

	responderKey, err := security.DeriveKey(responderSecret, salt, info)
	if err != nil {
		t.Fatalf("Responder DeriveKey failed: %v", err)
	}

	if len(initiatorKey) != 32 {
		t.Fatalf("Expected 32-byte AES key, got %d bytes", len(initiatorKey))
	}

	if !bytes.Equal(initiatorKey, responderKey) {
		t.Fatalf("HKDF derived keys do not match!")
	}
}

func TestAESGCMEncryptDecrypt(t *testing.T) {
	initiatorKeyPair, _ := security.GenerateKeyPair()
	responderKeyPair, _ := security.GenerateKeyPair()

	secret, _ := security.ComputeSharedSecret(initiatorKeyPair.PrivKey, responderKeyPair.PubKeyBytes())
	key, _ := security.DeriveKey(secret, nil, nil)

	plaintext := []byte("Sensitive GoShare File Chunk Content #1042")
	aad := []byte("chunk-header-id-1042")

	// Encrypt
	ciphertext, err := security.Encrypt(key, plaintext, aad)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := security.Decrypt(key, ciphertext, aad)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text = %s; want %s", string(decrypted), string(plaintext))
	}
}

func TestTamperDetection(t *testing.T) {
	initiatorKeyPair, _ := security.GenerateKeyPair()
	responderKeyPair, _ := security.GenerateKeyPair()

	secret, _ := security.ComputeSharedSecret(initiatorKeyPair.PrivKey, responderKeyPair.PubKeyBytes())
	key, _ := security.DeriveKey(secret, nil, nil)

	plaintext := []byte("Confidential data payload")

	ciphertext, err := security.Encrypt(key, plaintext, nil)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Tamper with the last byte of ciphertext
	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err = security.Decrypt(key, ciphertext, nil)
	if err == nil {
		t.Fatalf("Expected Decrypt to fail on tampered ciphertext, but it succeeded!")
	}

	if err != security.ErrAuthenticationFailed {
		t.Errorf("Expected ErrAuthenticationFailed, got: %v", err)
	}
}

func TestPassphraseGeneration(t *testing.T) {
	passphrase, err := security.GeneratePassphrase()
	if err != nil {
		t.Fatalf("GeneratePassphrase failed: %v", err)
	}

	words := strings.Split(passphrase, "-")
	if len(words) != 4 {
		t.Fatalf("Expected 4 words in passphrase, got %d: %s", len(words), passphrase)
	}

	if !security.VerifyPassphrase(passphrase, passphrase) {
		t.Errorf("VerifyPassphrase failed for identical passphrase")
	}

	if !security.VerifyPassphrase(strings.ToUpper(passphrase), passphrase) {
		t.Errorf("VerifyPassphrase failed case-insensitivity test")
	}
}

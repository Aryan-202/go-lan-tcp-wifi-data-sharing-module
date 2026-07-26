package security_test

import (
	"bytes"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

func TestGenerateKeyPairAndPubKeyBytes(t *testing.T) {
	kp, err := security.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if kp == nil || kp.PrivKey == nil || kp.PubKey == nil {
		t.Fatalf("Expected non-nil KeyPair, PrivKey, and PubKey")
	}

	pubBytes := kp.PubKeyBytes()
	if len(pubBytes) != 32 {
		t.Fatalf("Expected 32-byte public key, got %d bytes", len(pubBytes))
	}
}

func TestComputeSharedSecret(t *testing.T) {
	localKp, err := security.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair local failed: %v", err)
	}

	remoteKp, err := security.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair remote failed: %v", err)
	}

	secret1, err := security.ComputeSharedSecret(localKp.PrivKey, remoteKp.PubKeyBytes())
	if err != nil {
		t.Fatalf("ComputeSharedSecret secret1 failed: %v", err)
	}

	secret2, err := security.ComputeSharedSecret(remoteKp.PrivKey, localKp.PubKeyBytes())
	if err != nil {
		t.Fatalf("ComputeSharedSecret secret2 failed: %v", err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("Shared secrets do not match")
	}
}

func TestComputeSharedSecret_InvalidPubKeyBytes(t *testing.T) {
	localKp, _ := security.GenerateKeyPair()

	invalidBytes := []byte("too-short")
	_, err := security.ComputeSharedSecret(localKp.PrivKey, invalidBytes)
	if err == nil {
		t.Errorf("Expected error when passing invalid public key bytes, got nil")
	}
}

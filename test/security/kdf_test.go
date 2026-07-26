package security_test

import (
	"bytes"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

func TestDeriveKey_Success(t *testing.T) {
	secret := []byte("01234567890123456789012345678901")
	salt := []byte("custom-salt")
	info := []byte("custom-info")

	key1, err := security.DeriveKey(secret, salt, info)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if len(key1) != security.DefaultAESKeySize {
		t.Fatalf("Expected %d bytes key, got %d", security.DefaultAESKeySize, len(key1))
	}

	key2, err := security.DeriveKey(secret, salt, info)
	if err != nil {
		t.Fatalf("DeriveKey key2 failed: %v", err)
	}

	if !bytes.Equal(key1, key2) {
		t.Errorf("DeriveKey outputs are not deterministic for same inputs")
	}
}

func TestDeriveKey_DefaultSaltAndInfo(t *testing.T) {
	secret := []byte("01234567890123456789012345678901")

	key, err := security.DeriveKey(secret, nil, nil)
	if err != nil {
		t.Fatalf("DeriveKey with default salt/info failed: %v", err)
	}

	if len(key) != security.DefaultAESKeySize {
		t.Fatalf("Expected %d bytes key, got %d", security.DefaultAESKeySize, len(key))
	}
}

func TestDeriveKey_EmptySharedSecretError(t *testing.T) {
	_, err := security.DeriveKey(nil, nil, nil)
	if err == nil {
		t.Errorf("Expected error for empty shared secret, got nil")
	}
}

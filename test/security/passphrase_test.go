package security_test

import (
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

func TestGeneratePassphrase(t *testing.T) {
	passphrase, err := security.GeneratePassphrase()
	if err != nil {
		t.Fatalf("GeneratePassphrase failed: %v", err)
	}

	parts := strings.Split(passphrase, "-")
	if len(parts) != 4 {
		t.Fatalf("Expected 4 hyphen-separated words, got %d from: %s", len(parts), passphrase)
	}

	for _, word := range parts {
		if len(word) == 0 {
			t.Errorf("Empty word found in passphrase: %s", passphrase)
		}
	}
}

func TestVerifyPassphrase(t *testing.T) {
	expected := "orbit-falcon-amber-crest"

	// Exact match
	if !security.VerifyPassphrase("orbit-falcon-amber-crest", expected) {
		t.Errorf("Exact match failed")
	}

	// Case insensitivity
	if !security.VerifyPassphrase("ORBIT-FALCON-AMBER-CREST", expected) {
		t.Errorf("Case insensitive match failed")
	}

	// Whitespace trimming
	if !security.VerifyPassphrase("  orbit-falcon-amber-crest \n", expected) {
		t.Errorf("Whitespace trimming match failed")
	}

	// Mismatched passphrase
	if security.VerifyPassphrase("wrong-falcon-amber-crest", expected) {
		t.Errorf("Expected false for mismatched passphrase, got true")
	}
}

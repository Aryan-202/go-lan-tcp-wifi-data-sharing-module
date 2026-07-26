package discovery_test

import (
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
)

func TestStartAnnouncerAndStop(t *testing.T) {
	announcer, err := discovery.StartAnnouncer(8829)
	if err != nil {
		t.Fatalf("StartAnnouncer failed: %v", err)
	}

	if announcer == nil {
		t.Fatalf("Expected non-nil Announcer")
	}

	// Safe stop
	announcer.Stop()
}

func TestNilAnnouncerStop(t *testing.T) {
	var nilAnnouncer *discovery.Announcer
	// Must not panic
	nilAnnouncer.Stop()
}

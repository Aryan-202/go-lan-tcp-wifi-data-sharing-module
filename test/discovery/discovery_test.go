package discovery_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
)

func TestAnnounceAndDiscover(t *testing.T) {
	announcer, err := discovery.StartAnnouncer(8829)
	if err != nil {
		t.Fatalf("StartAnnouncer failed: %v", err)
	}
	defer announcer.Stop()

	ctx := context.Background()
	peers, err := discovery.DiscoverPeers(ctx, 2*time.Second)
	if err != nil {
		t.Fatalf("DiscoverPeers failed: %v", err)
	}

	// Should discover at least self or zero peers if loopback multicast is filtered by OS
	t.Logf("Discovered %d peer(s)", len(peers))
	for _, p := range peers {
		if p.Name == "" {
			t.Errorf("Expected non-empty peer name, got empty")
		}
		if !strings.Contains(p.String(), p.Name) {
			t.Errorf("Peer.String() = %q; expected to contain name %q", p.String(), p.Name)
		}
	}
}

func TestPeerFormatting(t *testing.T) {
	peer := discovery.Peer{
		ID:       "test-id",
		Name:     "TestDevice",
		IP:       "192.168.1.100",
		Port:     8829,
		OS:       "linux",
		Metadata: map[string]string{"version": "v1.0.0"},
	}

	formatted := peer.String()
	if !strings.Contains(formatted, "TestDevice") || !strings.Contains(formatted, "192.168.1.100:8829") {
		t.Errorf("Unexpected formatting: %s", formatted)
	}
}

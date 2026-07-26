package discovery_test

import (
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
)

func TestPeerString_WithMetadata(t *testing.T) {
	peer := discovery.Peer{
		ID:       "peer-1",
		Name:     "ArchLaptop",
		IP:       "192.168.1.50",
		Port:     8829,
		OS:       "linux",
		Metadata: map[string]string{"version": "v1.0.0"},
	}

	str := peer.String()
	if !strings.Contains(str, "ArchLaptop") || !strings.Contains(str, "192.168.1.50:8829") || !strings.Contains(str, "version=v1.0.0") {
		t.Errorf("Unexpected String() output: %s", str)
	}
}

func TestPeerString_WithoutMetadata(t *testing.T) {
	peer := discovery.Peer{
		ID:   "peer-2",
		Name: "MacBookPro",
		IP:   "192.168.1.51",
		Port: 8829,
		OS:   "darwin",
	}

	str := peer.String()
	if !strings.Contains(str, "MacBookPro [darwin] @ 192.168.1.51:8829") {
		t.Errorf("Unexpected String() output: %s", str)
	}
}

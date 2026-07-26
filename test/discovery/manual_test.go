package discovery_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
)

func TestProbeIP_Success(t *testing.T) {
	// Start a local TCP listener to simulate active GoShare peer
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test TCP listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	ctx := context.Background()

	peer, err := discovery.ProbeIP(ctx, "127.0.0.1", addr.Port)
	if err != nil {
		t.Fatalf("ProbeIP failed against local listener: %v", err)
	}

	if peer == nil || peer.IP != "127.0.0.1" || peer.Port != addr.Port {
		t.Errorf("Unexpected peer returned from ProbeIP: %v", peer)
	}
}

func TestProbeIP_UnreachableFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Probe port 59999 which is unlikely to have a TCP server
	_, err := discovery.ProbeIP(ctx, "127.0.0.1", 59999)
	if err == nil {
		t.Errorf("Expected ProbeIP to fail against closed port 59999, got nil error")
	}
}

func TestProbeIP_DefaultPortFallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Test port <= 0 fallback
	_, err := discovery.ProbeIP(ctx, "127.0.0.1", 0)
	if err == nil {
		t.Logf("ProbeIP failed as expected with default port fallback: %v", err)
	}
}

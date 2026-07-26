package discovery

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ProbeIP attempts a manual TCP probe to a given target IP address and port
func ProbeIP(ctx context.Context, targetIP string, port int) (*Peer, error) {
	if port <= 0 {
		port = DefaultPort
	}

	address := fmt.Sprintf("%s:%d", targetIP, port)
	dialer := net.Dialer{Timeout: 2 * time.Second}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("manual IP probe failed to reach %s: %w", address, err)
	}
	_ = conn.Close()

	return &Peer{
		ID:       fmt.Sprintf("manual-%s", targetIP),
		Name:     fmt.Sprintf("Peer@%s", targetIP),
		IP:       targetIP,
		Port:     port,
		OS:       "Manual Entry",
		Metadata: map[string]string{"type": "manual"},
	}, nil
}

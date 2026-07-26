package discovery

import (
	"fmt"
	"os"
	"runtime"

	"github.com/grandcat/zeroconf"
)

// Announcer manages mDNS service registration for GoShare
type Announcer struct {
	server *zeroconf.Server
}

// StartAnnouncer registers the GoShare mDNS service on local network
func StartAnnouncer(port int) (*Announcer, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "GoShare-Peer"
	}

	txtRecords := []string{
		fmt.Sprintf("os=%s", runtime.GOOS),
		"version=v1.0.0",
		"app=goshare",
	}

	server, err := zeroconf.Register(
		hostname,
		ServiceType,
		Domain,
		port,
		txtRecords,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register mDNS zeroconf service: %w", err)
	}

	return &Announcer{server: server}, nil
}

// Stop shuts down the mDNS service advertisement
func (a *Announcer) Stop() {
	if a != nil && a.server != nil {
		a.server.Shutdown()
	}
}

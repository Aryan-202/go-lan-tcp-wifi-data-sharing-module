package discovery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

// DiscoverPeers searches the local network for active GoShare peers using mDNS
func DiscoverPeers(ctx context.Context, timeout time.Duration) ([]Peer, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mDNS resolver: %w", err)
	}

	browseCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	entries := make(chan *zeroconf.ServiceEntry)
	discoveredMap := make(map[string]Peer)

	go func() {
		for entry := range entries {
			peer := parseServiceEntry(entry)
			if peer.IP != "" {
				key := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
				discoveredMap[key] = peer
			}
		}
	}()

	err = resolver.Browse(browseCtx, ServiceType, Domain, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to browse mDNS services: %w", err)
	}

	<-browseCtx.Done()

	peers := make([]Peer, 0, len(discoveredMap))
	for _, p := range discoveredMap {
		peers = append(peers, p)
	}

	return peers, nil
}

func parseServiceEntry(entry *zeroconf.ServiceEntry) Peer {
	peer := Peer{
		ID:       entry.Instance,
		Name:     entry.Instance,
		Port:     entry.Port,
		Metadata: make(map[string]string),
	}

	// Prefer IPv4 address
	if len(entry.AddrIPv4) > 0 {
		peer.IP = entry.AddrIPv4[0].String()
	} else if len(entry.AddrIPv6) > 0 {
		peer.IP = entry.AddrIPv6[0].String()
	}

	for _, txt := range entry.Text {
		parts := strings.SplitN(txt, "=", 2)
		if len(parts) == 2 {
			key, val := parts[0], parts[1]
			if key == "os" {
				peer.OS = val
			} else {
				peer.Metadata[key] = val
			}
		}
	}

	if peer.OS == "" {
		peer.OS = "Unknown OS"
	}

	return peer
}

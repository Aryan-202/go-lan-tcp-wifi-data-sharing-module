package discovery

import (
	"fmt"
	"strings"
)

// ServiceType is the mDNS service type used by GoShare
const ServiceType = "_goshare._tcp"

// Domain is the mDNS domain
const Domain = "local."

// DefaultPort is the default TCP/UDP port for GoShare communication
const DefaultPort = 8829

// Peer represents a discovered GoShare instance on the local network
type Peer struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	IP       string            `json:"ip"`
	Port     int               `json:"port"`
	OS       string            `json:"os"`
	Metadata map[string]string `json:"metadata"`
}

// String returns a human-readable representation of a Peer
func (p Peer) String() string {
	var metaStr string
	if len(p.Metadata) > 0 {
		var pairs []string
		for k, v := range p.Metadata {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
		metaStr = fmt.Sprintf(" (%s)", strings.Join(pairs, ", "))
	}
	return fmt.Sprintf("%s [%s] @ %s:%d%s", p.Name, p.OS, p.IP, p.Port, metaStr)
}

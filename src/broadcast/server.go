package broadcast

import (
	"fmt"
	"log"
	"net"
)

func DiscoverServer() error {
	localAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:8829")
	if err != nil {
		return fmt.Errorf("failed to resolve local address: %w", err)
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return fmt.Errorf("failed to bind to port: %w", err)
	}
	defer conn.Close()

	log.Printf("Listening for UDP broadcasts on port %s...", localAddr.String())

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		log.Printf("Received %d bytes from %s: %s", n, remoteAddr, string(buffer[:n]))
	}
}
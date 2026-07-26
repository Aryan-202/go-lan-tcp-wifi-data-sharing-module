package discovery

import (
	"fmt"
	"log"
	"net"
	"time"
)

func DiscoverClient() error {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8829")
	if err != nil {
		return fmt.Errorf("failed to resolve broadcast address: %w", err)
	}

	conn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		return fmt.Errorf("failed to dial udp: %w", err)
	}
	defer conn.Close()

	log.Println("Starting UDP broadcast loop...")
	message := []byte("Hello from GoShare broadcast!")

	for {
		_, err := conn.Write(message)
		if err != nil {
			log.Printf("Broadcast write error: %v", err)
		} else {
			log.Println("Broadcast message sent successfully")
		}

		time.Sleep(3 * time.Second)
	}
}


package discovery

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"go-lan-tcp-sharing/pkg/models"
)

const (
	MagicNumber = 0x50325078
	Port        = 9998
)

type Message struct {
	Magic   uint32      `json:"magic"`
	Type    string      `json:"type"`
	Version string      `json:"version"`
	Peer    models.Peer `json:"peer"`
}

type Service struct {
	peers  map[string]models.Peer
	mutex  sync.RWMutex
	myPeer models.Peer
	stop   chan struct{}
}

func NewService(myID string, myPort int) *Service {
	return &Service{
		peers: make(map[string]models.Peer),
		myPeer: models.Peer{
			ID:   myID,
			Port: myPort,
		},
		stop: make(chan struct{}),
	}
}

func (s *Service) Start() error {
	addr := &net.UDPAddr{
		Port: Port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	go s.listen(conn)
	go s.announce()
	go s.cleanup()

	return nil
}

func (s *Service) Stop() {
	close(s.stop)
}

func (s *Service) announce() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	baddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: Port,
	}
	conn, err := net.DialUDP("udp", nil, baddr)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		select {
		case <-s.stop:
			s.myPeer.IP = GetLocalIP()
			msg := Message{
				Magic:   MagicNumber,
				Type:    "GOODBYE",
				Version: "1.0",
				Peer:    s.myPeer,
			}
			b, _ := json.Marshal(msg)
			conn.Write(b)
			return
		case <-ticker.C:
			s.myPeer.IP = GetLocalIP()
			msg := Message{
				Magic:   MagicNumber,
				Type:    "ANNOUNCE",
				Version: "1.0",
				Peer:    s.myPeer,
			}
			b, _ := json.Marshal(msg)
			conn.Write(b)
		}
	}
}

func (s *Service) listen(conn *net.UDPConn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		select {
		case <-s.stop:
			return
		default:
		}
		
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var msg Message
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			continue
		}
		if msg.Magic != MagicNumber {
			continue
		}

		if msg.Type == "ANNOUNCE" || msg.Type == "RESPONSE" {
			if msg.Peer.ID == s.myPeer.ID {
				continue
			}
			msg.Peer.IP = addr.IP.String()
			msg.Peer.LastSeen = time.Now()

			s.mutex.Lock()
			s.peers[msg.Peer.ID] = msg.Peer
			s.mutex.Unlock()
		} else if msg.Type == "GOODBYE" {
			s.mutex.Lock()
			delete(s.peers, msg.Peer.ID)
			s.mutex.Unlock()
		} else if msg.Type == "QUERY" {
			if msg.Peer.ID != s.myPeer.ID {
				s.myPeer.IP = GetLocalIP()
				respMsg := Message{
					Magic:   MagicNumber,
					Type:    "RESPONSE",
					Version: "1.0",
					Peer:    s.myPeer,
				}
				b, _ := json.Marshal(respMsg)
				
				uaddr := &net.UDPAddr{
					IP:   addr.IP,
					Port: Port,
				}
				
				respConn, _ := net.DialUDP("udp", nil, uaddr)
				if respConn != nil {
					respConn.Write(b)
					respConn.Close()
				}
			}
		}
	}
}

func (s *Service) cleanup() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			now := time.Now()
			s.mutex.Lock()
			for id, peer := range s.peers {
				if now.Sub(peer.LastSeen) > 15*time.Second {
					delete(s.peers, id)
				}
			}
			s.mutex.Unlock()
		}
	}
}

func (s *Service) GetPeers() []models.Peer {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	peers := make([]models.Peer, 0, len(s.peers))
	for _, p := range s.peers {
		peers = append(peers, p)
	}
	return peers
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

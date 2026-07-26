package transfer

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

// SendFile encrypts and transmits a local file over a TCP connection
func SendFile(conn net.Conn, filePath string, sessionKey []byte, state *TransferState) (*FileManifest, error) {
	defer conn.Close()

	// 1. Generate Manifest
	manifest, err := GenerateManifest(filePath, DefaultSegmentSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate manifest for send: %w", err)
	}

	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// 2. Encrypt & Send Manifest
	encryptedManifest, err := security.Encrypt(sessionKey, manifestBytes, []byte("manifest"))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt manifest: %w", err)
	}

	if err := WritePacket(conn, encryptedManifest); err != nil {
		return nil, fmt.Errorf("failed to transmit manifest packet: %w", err)
	}

	// 3. Open source file for reading segments
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file for transmission: %w", err)
	}
	defer file.Close()

	// 4. Send segments
	for i := 0; i < manifest.TotalSegments; i++ {
		if state != nil && state.IsSegmentCompleted(i) {
			continue // Skip already transferred segment when resuming
		}

		segmentData, _, err := ReadSegment(file, i, manifest.SegmentSize)
		if err != nil {
			return nil, fmt.Errorf("failed to read segment %d: %w", i, err)
		}

		aad := fmt.Appendf(nil, "segment-%d", i)
		encryptedPayload, err := security.Encrypt(sessionKey, segmentData, aad)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt segment %d: %w", i, err)
		}

		packet := SegmentPacket{
			FileID:       manifest.FileID,
			SegmentIndex: i,
			DataLength:   len(segmentData),
			Payload:      encryptedPayload,
		}

		packetBytes, err := json.Marshal(packet)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal segment packet %d: %w", i, err)
		}

		if err := WritePacket(conn, packetBytes); err != nil {
			return nil, fmt.Errorf("failed to transmit segment packet %d: %w", i, err)
		}
	}

	return manifest, nil
}

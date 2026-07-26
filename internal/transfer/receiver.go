package transfer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
)

// ReceiveFile accepts an incoming encrypted file transfer over a TCP connection
func ReceiveFile(conn net.Conn, saveDir string, sessionKey []byte, state *TransferState) (*FileManifest, error) {
	defer conn.Close()

	// 1. Read encrypted Manifest packet
	encryptedManifest, err := ReadPacket(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to receive file manifest packet: %w", err)
	}

	manifestBytes, err := security.Decrypt(sessionKey, encryptedManifest, []byte("manifest"))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file manifest: %w", err)
	}

	var manifest FileManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file manifest: %w", err)
	}

	// 2. Prepare target output file on disk
	savePath := filepath.Join(saveDir, manifest.FileName)
	outFile, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer outFile.Close()

	if err := outFile.Truncate(manifest.FileSize); err != nil {
		return nil, fmt.Errorf("failed to pre-allocate destination file size: %w", err)
	}

	// 3. Receive segment packets
	receivedCount := 0
	if state == nil {
		state = NewTransferState(manifest.FileID, manifest.FileName)
	}

	for receivedCount < manifest.TotalSegments {
		packetBytes, err := ReadPacket(conn)
		if err != nil {
			if err == io.EOF && receivedCount == manifest.TotalSegments {
				break
			}
			return nil, fmt.Errorf("failed to read segment packet: %w", err)
		}

		var packet SegmentPacket
		if err := json.Unmarshal(packetBytes, &packet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal segment packet: %w", err)
		}

		// Decrypt segment payload with segment AAD header
		aad := fmt.Appendf([]byte{}, "segment-%d", packet.SegmentIndex)
		decryptedSegment, err := security.Decrypt(sessionKey, packet.Payload, aad)
		if err != nil {
			return nil, fmt.Errorf("decryption failed for segment index %d: %w", packet.SegmentIndex, err)
		}

		// Check SHA-256 segment integrity
		segmentHash := sha256.Sum256(decryptedSegment)
		hashStr := hex.EncodeToString(segmentHash[:])
		if hashStr != manifest.SegmentHashes[packet.SegmentIndex] {
			return nil, fmt.Errorf("integrity check failed for segment %d: expected %s, got %s",
				packet.SegmentIndex, manifest.SegmentHashes[packet.SegmentIndex], hashStr)
		}

		// Write segment to thread-safe offset
		offset := int64(packet.SegmentIndex) * manifest.SegmentSize
		if _, err := outFile.WriteAt(decryptedSegment, offset); err != nil {
			return nil, fmt.Errorf("failed to write segment %d to disk offset %d: %w", packet.SegmentIndex, offset, err)
		}

		state.MarkSegmentCompleted(packet.SegmentIndex)
		receivedCount++
	}

	return &manifest, nil
}

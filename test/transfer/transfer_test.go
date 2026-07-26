package transfer_test

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/transfer"
)

func TestGenerateManifestAndReadSegment(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample_test.dat")

	// 100 KB test data
	testData := make([]byte, 100*1024)
	_, _ = rand.Read(testData)
	if err := os.WriteFile(filePath, testData, 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	segmentSize := int64(16 * 1024) // 16 KB segments
	manifest, err := transfer.GenerateManifest(filePath, segmentSize)
	if err != nil {
		t.Fatalf("GenerateManifest failed: %v", err)
	}

	if manifest.TotalSegments != 7 { // ceil(100/16) = 7
		t.Fatalf("Expected 7 segments, got %d", manifest.TotalSegments)
	}

	expectedHash := sha256.Sum256(testData)
	if manifest.FileHash != hex.EncodeToString(expectedHash[:]) {
		t.Errorf("FileHash mismatch: expected %s, got %s", hex.EncodeToString(expectedHash[:]), manifest.FileHash)
	}

	// Verify reading a segment
	file, _ := os.Open(filePath)
	defer file.Close()

	segmentData, hash, err := transfer.ReadSegment(file, 0, segmentSize)
	if err != nil {
		t.Fatalf("ReadSegment failed: %v", err)
	}

	if len(segmentData) != int(segmentSize) {
		t.Errorf("Expected segment size %d, got %d", segmentSize, len(segmentData))
	}

	if hash != manifest.SegmentHashes[0] {
		t.Errorf("Segment 0 hash mismatch")
	}
}

func TestEncryptedFileTransfer(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "large_source.bin")
	recvDir := filepath.Join(tempDir, "received")
	_ = os.MkdirAll(recvDir, 0755)

	// 5 MB random test file
	fileSize := 5 * 1024 * 1024
	sourceData := make([]byte, fileSize)
	_, _ = rand.Read(sourceData)
	if err := os.WriteFile(sourcePath, sourceData, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Negotiate X25519 + HKDF session key
	kp1, _ := security.GenerateKeyPair()
	kp2, _ := security.GenerateKeyPair()
	secret, _ := security.ComputeSharedSecret(kp1.PrivKey, kp2.PubKeyBytes())
	sessionKey, _ := security.DeriveKey(secret, nil, nil)

	// Start local TCP listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen on TCP: %v", err)
	}
	defer listener.Close()

	errChan := make(chan error, 1)

	// Receiver goroutine
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		_, err = transfer.ReceiveFile(conn, recvDir, sessionKey, nil)
		errChan <- err
	}()

	// Sender client
	senderConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial receiver TCP: %v", err)
	}

	_, err = transfer.SendFile(senderConn, sourcePath, sessionKey, nil)
	if err != nil {
		t.Fatalf("SendFile failed: %v", err)
	}

	if recvErr := <-errChan; recvErr != nil {
		t.Fatalf("ReceiveFile failed: %v", recvErr)
	}

	// Verify received file matches source byte-for-byte
	receivedPath := filepath.Join(recvDir, "large_source.bin")
	receivedData, err := os.ReadFile(receivedPath)
	if err != nil {
		t.Fatalf("Failed to read received file: %v", err)
	}

	if !bytes.Equal(sourceData, receivedData) {
		t.Fatalf("Received file content does not match original file!")
	}
}

func TestTransferStatePersistence(t *testing.T) {
	tempDir := t.TempDir()
	statePath := filepath.Join(tempDir, "transfer_state.json")

	state := transfer.NewTransferState("file-id-999", "sample.zip")
	state.MarkSegmentCompleted(0)
	state.MarkSegmentCompleted(1)
	state.MarkSegmentCompleted(4)

	if err := transfer.SaveState(statePath, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	loadedState, err := transfer.LoadState(statePath)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if loadedState.FileID != "file-id-999" || loadedState.FileName != "sample.zip" {
		t.Errorf("Unexpected loaded state header: %v", loadedState)
	}

	if !loadedState.IsSegmentCompleted(0) || !loadedState.IsSegmentCompleted(1) || !loadedState.IsSegmentCompleted(4) {
		t.Errorf("Expected completed segments 0, 1, 4 to be true")
	}

	if loadedState.IsSegmentCompleted(2) {
		t.Errorf("Expected completed segment 2 to be false")
	}
}

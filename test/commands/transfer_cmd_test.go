package commands_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/cmd"
)

func TestSendAndReceiveCommands(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "cli_test_data.txt")
	downloadDir := filepath.Join(tempDir, "downloads")

	// Create test file
	testContent := make([]byte, 50*1024)
	_, _ = rand.Read(testContent)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	receiveCmd := cmd.GetReceiveCmd()
	sendCmd := cmd.GetSendCmd()

	_ = receiveCmd.Flags().Set("port", "18829")
	_ = receiveCmd.Flags().Set("dir", downloadDir)

	_ = sendCmd.Flags().Set("ip", "127.0.0.1")
	_ = sendCmd.Flags().Set("port", "18829")

	recvBuf := new(bytes.Buffer)
	sendBuf := new(bytes.Buffer)

	receiveCmd.SetOut(recvBuf)
	receiveCmd.SetErr(recvBuf)
	sendCmd.SetOut(sendBuf)
	sendCmd.SetErr(sendBuf)

	recvErrChan := make(chan error, 1)

	// Run receiver in background goroutine
	go func() {
		err := receiveCmd.RunE(receiveCmd, []string{})
		recvErrChan <- err
	}()

	// Give receiver time to bind port
	time.Sleep(200 * time.Millisecond)

	// Run sender
	err := sendCmd.RunE(sendCmd, []string{sourceFile})
	if err != nil {
		t.Fatalf("sendCmd.RunE failed: %v", err)
	}

	if recvErr := <-recvErrChan; recvErr != nil {
		t.Fatalf("receiveCmd.RunE failed: %v", recvErr)
	}

	// Verify file received
	receivedPath := filepath.Join(downloadDir, "cli_test_data.txt")
	receivedBytes, err := os.ReadFile(receivedPath)
	if err != nil {
		t.Fatalf("Failed to read received file: %v", err)
	}

	if !bytes.Equal(testContent, receivedBytes) {
		t.Fatalf("Received file content does not match source content!")
	}
}

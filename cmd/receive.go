package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/transfer"
	"github.com/spf13/cobra"
)

var (
	receivePortFlag int
	receiveDirFlag  string
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Start GoShare receiver to accept files from peers",
	Long:  `Start listening for incoming GoShare file transfers over TCP and mDNS discovery.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if receiveDirFlag == "" {
			receiveDirFlag = "./downloads"
		}
		if err := os.MkdirAll(receiveDirFlag, 0755); err != nil {
			return fmt.Errorf("failed to create download directory: %w", err)
		}

		// 1. Start mDNS advertiser
		announcer, err := discovery.StartAnnouncer(receivePortFlag)
		if err != nil {
			cmd.Printf("Warning: Failed to start mDNS advertiser: %v\n", err)
		} else {
			defer announcer.Stop()
			cmd.Println("mDNS Service active: Discoverable by nearby peers")
		}

		// 2. Generate key pair and passphrase
		localKeyPair, err := security.GenerateKeyPair()
		if err != nil {
			return fmt.Errorf("failed to generate security keys: %w", err)
		}

		passphrase, err := security.GeneratePassphrase()
		if err != nil {
			return fmt.Errorf("failed to generate passphrase: %w", err)
		}

		// 3. Start TCP listener
		addrStr := fmt.Sprintf("0.0.0.0:%d", receivePortFlag)
		listener, err := net.Listen("tcp", addrStr)
		if err != nil {
			return fmt.Errorf("failed to listen on TCP %s: %w", addrStr, err)
		}
		defer listener.Close()

		cmd.Println("==================================================")
		cmd.Printf("  GoShare Receiver Listening on Port %d\n", receivePortFlag)
		cmd.Printf("  Save Directory    : %s\n", receiveDirFlag)
		cmd.Printf("  Pairing Passphrase: %s\n", passphrase)
		cmd.Println("==================================================")
		cmd.Println("Waiting for incoming connection...")

		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept peer connection: %w", err)
		}

		cmd.Printf("Peer connected from %s! Performing key exchange...\n", conn.RemoteAddr().String())

		// Perform X25519 Key Exchange over socket
		// Sender sends public key bytes (32 bytes)
		peerPubBytes := make([]byte, 32)
		if _, err := io.ReadFull(conn, peerPubBytes); err != nil {
			return fmt.Errorf("failed to read sender public key: %w", err)
		}

		// Send receiver public key bytes (32 bytes)
		if _, err := conn.Write(localKeyPair.PubKeyBytes()); err != nil {
			return fmt.Errorf("failed to send receiver public key: %w", err)
		}

		// Derive session key
		rawSecret, err := security.ComputeSharedSecret(localKeyPair.PrivKey, peerPubBytes)
		if err != nil {
			return fmt.Errorf("failed to compute shared secret: %w", err)
		}

		sessionKey, err := security.DeriveKey(rawSecret, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to derive session key: %w", err)
		}

		cmd.Println("Key exchange complete. Receiving encrypted file...")
		manifest, err := transfer.ReceiveFile(conn, receiveDirFlag, sessionKey, nil)
		if err != nil {
			return fmt.Errorf("file transfer failed: %w", err)
		}

		cmd.Println("==================================================")
		cmd.Printf("  SUCCESS: Received '%s' (%d bytes)\n", manifest.FileName, manifest.FileSize)
		cmd.Printf("  Saved to: %s\n", filepath.Join(receiveDirFlag, manifest.FileName))
		cmd.Println("==================================================")

		return nil
	},
}

func init() {
	receiveCmd.Flags().IntVarP(&receivePortFlag, "port", "p", 8829, "TCP port to listen on")
	receiveCmd.Flags().StringVarP(&receiveDirFlag, "dir", "d", "./downloads", "Directory to save received files")
	rootCmd.AddCommand(receiveCmd)
}

func GetReceiveCmd() *cobra.Command {
	return receiveCmd
}

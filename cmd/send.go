package cmd

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/security"
	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/transfer"
	"github.com/spf13/cobra"
)

var (
	sendIPFlag         string
	sendPortFlag       int
	sendPassphraseFlag string
)

var sendCmd = &cobra.Command{
	Use:   "send <file_path>",
	Short: "Send a file to a GoShare peer",
	Long:  `Encrypt and transmit a file to a receiving GoShare peer over local network.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		if _, err := os.Stat(filePath); err != nil {
			return fmt.Errorf("invalid file path: %w", err)
		}

		if sendIPFlag == "" {
			return fmt.Errorf("destination IP address is required (use --ip <target_ip>)")
		}

		// 1. Generate key pair
		localKeyPair, err := security.GenerateKeyPair()
		if err != nil {
			return fmt.Errorf("failed to generate key pair: %w", err)
		}

		// 2. Connect to receiver over TCP
		targetAddr := fmt.Sprintf("%s:%d", sendIPFlag, sendPortFlag)
		cmd.Printf("Connecting to GoShare peer at %s...\n", targetAddr)

		conn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			return fmt.Errorf("failed to connect to peer at %s: %w", targetAddr, err)
		}

		// Perform X25519 Key Exchange
		// Send sender public key bytes (32 bytes)
		if _, err := conn.Write(localKeyPair.PubKeyBytes()); err != nil {
			conn.Close()
			return fmt.Errorf("failed to send public key: %w", err)
		}

		// Read receiver public key bytes (32 bytes)
		peerPubBytes := make([]byte, 32)
		if _, err := io.ReadFull(conn, peerPubBytes); err != nil {
			conn.Close()
			return fmt.Errorf("failed to read receiver public key: %w", err)
		}

		// Derive session key
		rawSecret, err := security.ComputeSharedSecret(localKeyPair.PrivKey, peerPubBytes)
		if err != nil {
			conn.Close()
			return fmt.Errorf("failed to compute shared secret: %w", err)
		}

		sessionKey, err := security.DeriveKey(rawSecret, nil, nil)
		if err != nil {
			conn.Close()
			return fmt.Errorf("failed to derive session key: %w", err)
		}

		cmd.Println("Encrypted session established. Transferring file segments...")
		manifest, err := transfer.SendFile(conn, filePath, sessionKey, nil)
		if err != nil {
			return fmt.Errorf("failed to send file: %w", err)
		}

		cmd.Println("==================================================")
		cmd.Printf("  SUCCESS: Transmitted '%s' (%d bytes, %d segments)\n", manifest.FileName, manifest.FileSize, manifest.TotalSegments)
		cmd.Println("==================================================")

		return nil
	},
}

func init() {
	sendCmd.Flags().StringVar(&sendIPFlag, "ip", "", "Destination peer IP address")
	sendCmd.Flags().IntVarP(&sendPortFlag, "port", "p", 8829, "Destination peer TCP port")
	sendCmd.Flags().StringVarP(&sendPassphraseFlag, "passphrase", "s", "", "4-word pairing passphrase printed by receiver")
	rootCmd.AddCommand(sendCmd)
}

func GetSendCmd() *cobra.Command {
	return sendCmd
}

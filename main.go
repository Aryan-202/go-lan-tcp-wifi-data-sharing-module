package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-lan-tcp-sharing/pkg/discovery"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "p2pxfer",
	Short: "P2Pxfer is a CLI-first local file transfer system",
	Long:  `A decentralized, peer-to-peer file transfer utility for the local network.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	port        int
	chunkSize   int
	parallel    int
	encrypt     bool
	compress    bool
	outputDir   string
	peerAddress string
	authCode    string
)

func init() {
	// send command
	sendCmd := &cobra.Command{
		Use:   "send [file]",
		Short: "Send a file to a peer",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Sending file:", args[0])
			// TODO: Implement send
		},
	}
	sendCmd.Flags().IntVar(&port, "port", 9999, "Custom TCP port to listen on")
	sendCmd.Flags().IntVar(&chunkSize, "chunk-size", 4*1024*1024, "Override default 4MB chunk size")
	sendCmd.Flags().IntVar(&parallel, "parallel", 4, "Number of parallel transfer workers")
	sendCmd.Flags().BoolVar(&encrypt, "encrypt", true, "Enable/Disable encryption")
	sendCmd.Flags().BoolVar(&compress, "compress", true, "Enable/Disable compression (auto-detect)")
	rootCmd.AddCommand(sendCmd)

	// receive command
	receiveCmd := &cobra.Command{
		Use:   "receive",
		Short: "Receive a file from a peer",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Receiving file...")
			// TODO: Implement receive
		},
	}
	receiveCmd.Flags().IntVar(&port, "port", 9999, "Custom TCP port to listen on")
	receiveCmd.Flags().StringVar(&outputDir, "output", ".", "Output directory for received files")
	receiveCmd.Flags().StringVar(&peerAddress, "peer", "", "Manually specify a peer IP:Port")
	receiveCmd.Flags().StringVar(&authCode, "code", "", "Verification code for secure connection")
	rootCmd.AddCommand(receiveCmd)

	// discover command
	discoverCmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover peers on the local network",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Discovering peers...")
			myID := fmt.Sprintf("peer-%d", time.Now().UnixNano())
			svc := discovery.NewService(myID, port)
			if err := svc.Start(); err != nil {
				fmt.Println("Error starting discovery:", err)
				return
			}
			defer svc.Stop()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-c:
					return
				case <-ticker.C:
					peers := svc.GetPeers()
					fmt.Printf("\rFound %d peers. Press Ctrl+C to exit. ", len(peers))
					if len(peers) > 0 {
						fmt.Println("\n--- Peers ---")
						for _, p := range peers {
							fmt.Printf("- ID: %s, IP: %s:%d\n", p.ID, p.IP, p.Port)
						}
						fmt.Println("-------------")
					}
				}
			}
		},
	}
	rootCmd.AddCommand(discoverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

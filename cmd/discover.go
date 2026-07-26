package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/internal/discovery"
	"github.com/spf13/cobra"
)

var (
	timeoutFlag  time.Duration
	targetIPFlag string
	announceFlag bool
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover nearby GoShare devices",
	Long: `Search the local network for GoShare peers using mDNS
service discovery or manual IP probing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Searching for nearby devices...")

		ctx := context.Background()

		if announceFlag {
			announcer, err := discovery.StartAnnouncer(discovery.DefaultPort)
			if err != nil {
				cmd.Printf("Warning: Unable to start mDNS advertiser: %v\n", err)
			} else {
				defer announcer.Stop()
			}
		}

		if targetIPFlag != "" {
			cmd.Printf("Probing manual target IP: %s...\n", targetIPFlag)
			peer, err := discovery.ProbeIP(ctx, targetIPFlag, discovery.DefaultPort)
			if err != nil {
				return fmt.Errorf("manual IP probe error: %w", err)
			}
			cmd.Printf("Found peer: %s\n", peer.String())
			return nil
		}

		peers, err := discovery.DiscoverPeers(ctx, timeoutFlag)
		if err != nil {
			return fmt.Errorf("discovery failed: %w", err)
		}

		if len(peers) == 0 {
			cmd.Println("No nearby GoShare peers found.")
			return nil
		}

		cmd.Printf("Discovered %d peer(s):\n", len(peers))
		for i, p := range peers {
			cmd.Printf("  [%d] %s\n", i+1, p.String())
		}

		return nil
	},
}

func init() {
	discoverCmd.Flags().DurationVarP(&timeoutFlag, "timeout", "t", 5*time.Second, "Timeout duration for peer discovery")
	discoverCmd.Flags().StringVar(&targetIPFlag, "ip", "", "Target IP address for manual discovery fallback")
	discoverCmd.Flags().BoolVarP(&announceFlag, "announce", "a", true, "Announce self as a discoverable peer while searching")

	rootCmd.AddCommand(discoverCmd)
}

func GetDiscoverCmd() *cobra.Command {
	return discoverCmd
}
package cmd

import (
	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover nearby GoShare devices",
	Long: `Search the local network for GoShare peers using
the configured discovery protocols.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Searching for nearby devices...")

		// TODO: Call internal/discovery package

		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func GetDiscoverCmd() *cobra.Command {
	return discoverCmd
}
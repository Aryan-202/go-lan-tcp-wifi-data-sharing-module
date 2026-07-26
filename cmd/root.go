package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoShare",
	Short: "GoShare is a Go CLI tool that can share files at the maximum speed possible in a LAN",
	Long:  "GoShare is made for local area networks and is highly secure while sharing files across different OS environments like Mac, Windows, Linux, Android, iOS, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error while executing zero command '%s'\n", err)
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "goshare",
	Short:        "GoShare is a Go CLI tool that can share files at the maximum speed possible in a LAN",
	Long:         "GoShare is made for local area networks and is highly secure while sharing files across different OS environments like Mac, Windows, Linux, Android, iOS, etc.",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	// Sanitize os.Args in case shell/wrapper prepends executable path to argument list
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		if strings.Contains(firstArg, "/") || strings.HasPrefix(filepath.Base(firstArg), "goshare") {
			isCmd := false
			for _, c := range rootCmd.Commands() {
				if c.Name() == firstArg || c.Name() == filepath.Base(firstArg) {
					isCmd = true
					break
				}
			}
			if !isCmd && firstArg != "help" && firstArg != "completion" {
				os.Args = append(os.Args[:1], os.Args[2:]...)
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error while executing zero command '%s'\n", err)
		os.Exit(1)
	}
}

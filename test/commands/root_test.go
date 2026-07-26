package commands_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/cmd"
)

func TestRootCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	discoverCmd := cmd.GetDiscoverCmd()
	_ = discoverCmd

	// GetRootCmd executes root CLI
	rootCmd := cmd.GetDiscoverCmd().Root()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() --help failed: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "GoShare is made for local area networks") {
		t.Errorf("Unexpected rootCmd help output: %s", got)
	}
}

func TestExecuteHelper(t *testing.T) {
	// Execute() is the main entry point from main.go
	rootCmd := cmd.GetDiscoverCmd().Root()
	rootCmd.SetArgs([]string{"--help"})
	cmd.Execute()
}

package commands_test

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/cmd"
)

func TestDiscoverCmd(t *testing.T) {
	buf := new(bytes.Buffer)

	discoverCmd := cmd.GetDiscoverCmd()
	discoverCmd.SetOut(buf)
	discoverCmd.SetErr(buf)
	_ = discoverCmd.Flags().Set("timeout", "200ms")

	err := discoverCmd.RunE(discoverCmd, []string{})
	if err != nil {
		t.Fatalf("discoverCmd.RunE() failed unexpectedly: %v", err)
	}

	got := buf.String()
	want := "Searching for nearby devices..."

	if !strings.Contains(got, want) {
		t.Errorf("discoverCmd output = %q; want it to contain %q", got, want)
	}
}

func TestDiscoverCmd_ManualIP_Success(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	buf := new(bytes.Buffer)
	discoverCmd := cmd.GetDiscoverCmd()
	discoverCmd.SetOut(buf)
	discoverCmd.SetErr(buf)

	_ = discoverCmd.Flags().Set("ip", "127.0.0.1")
	_ = discoverCmd.Flags().Set("announce", "false")

	_ = discoverCmd.RunE(discoverCmd, []string{})

	got := buf.String()
	if !strings.Contains(got, "Probing manual target IP: 127.0.0.1...") {
		t.Errorf("Unexpected manual IP probe output: %s", got)
	}
}

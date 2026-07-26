package commands_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/cmd"
)

func TestDiscoverCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	
	// Fetch the command via the public helper function
	discoverCmd := cmd.GetDiscoverCmd()
	
	discoverCmd.SetOut(buf)
	discoverCmd.SetErr(buf)
	discoverCmd.SetArgs([]string{})

	err := discoverCmd.Execute()
	if err != nil {
		t.Fatalf("discoverCmd.Execute() failed unexpectedly: %v", err)
	}

	got := buf.String()
	want := "Searching for nearby devices..."

	if !strings.Contains(got, want) {
		t.Errorf("discoverCmd output = %q; want it to contain %q", got, want)
	}
}

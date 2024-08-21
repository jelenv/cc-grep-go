package main_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		input        string
		pattern      string
		expectedExit int
	}{
		{"apple", "a", 0},
		{"apple", "b", 1},
		// {"apple123", "/d", 1},
	}

	for _, tt := range tests {
		t.Run(tt.input+"-"+tt.pattern, func(t *testing.T) {
			cmd := exec.Command("../../build/mygrep", "-E", tt.pattern)

			cmd.Stdin = strings.NewReader(tt.input)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			exitCode := cmd.ProcessState.ExitCode()
			if err != nil && exitCode == -1 {
				t.Fatalf("error running command: %v", err)
			}

			if stdout.Len() > 0 {
				t.Logf("stdout:\n%s", stdout.String())
			}
			if stderr.Len() > 0 {
				t.Logf("stderr:\n%s", stderr.String())
			}

			if exitCode != tt.expectedExit {
				t.Fatalf("expected exit code %d but got %d", tt.expectedExit, exitCode)
			}
		})
	}
}

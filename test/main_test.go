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
		// 1. Match a literal character
		{"apple", "a", 0},
		{"apple", "b", 1},
		// 2. Match digits
		{"apple123", "\\d", 0},
		{"apple", "\\d", 1},
		// 3. Match alphanumeric characters
		{"alpha-num3ric", "\\w", 0},
		{"$!?", "\\w", 1},
		// 4. Positive character groups
		{"apple", "[aeiou]", 0},
		{"apple", "[xyz]", 1},
		// 5. Negative character groups
		{"apple", "[^xyz]", 0},
		{"apple", "[^aeiou]", 1},
		// 6. Combine character groups
		{"1 apple", "\\d apple", 0},
		{"sally has 1 dog", "\\d \\w\\w\\ws", 1},
		// 7. Start of string anchor
		{"log", "^log", 0},
		{"clog", "^log", 1},
		// 8. End of string anchor
		{"loglog", "log$", 0},
		{"logs", "log$", 1},
		// 9. Plus quantifier (1 or more)
		{"caats", "ca+ts", 0},
		{"log", "a+", 1},
		// 10. Question mark (0 or 1)
		{"dog", "dogs?", 0},
		{"cat", "ca?t", 0},
		{"cat", "dog?", 1},
	}

	for _, tt := range tests {
		t.Run(tt.input+"-"+tt.pattern, func(t *testing.T) {
			cmd := exec.Command("../build/mygrep", "-E", tt.pattern)

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

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestFlagParse(t *testing.T) {

	tests := []struct {
		args                     []string
		verbose, dryRun, dotFile bool
		path                     string
		exitCode                 int
	}{
		{
			args:     []string{"prog"},
			exitCode: 1,
		},
		{
			args:     []string{"prog", "a/path"},
			verbose:  false,
			dryRun:   false,
			dotFile:  false,
			path:     "a/path",
			exitCode: 0,
		},
		{
			args:     []string{"prog", "-v", "a/path"},
			verbose:  true,
			dryRun:   false,
			dotFile:  false,
			path:     "a/path",
			exitCode: 0,
		},
		{
			args:     []string{"prog", "-d", "a/path"},
			verbose:  false,
			dryRun:   true,
			dotFile:  false,
			path:     "a/path",
			exitCode: 0,
		},
		{
			args:     []string{"prog", "a/path", "another/path"},
			verbose:  false,
			dryRun:   false,
			dotFile:  false,
			path:     "",
			exitCode: 1, // second path unexpected
		},
		{
			args:     []string{"prog", "-d", "-v", "a/path"},
			verbose:  true,
			dryRun:   true,
			dotFile:  false,
			path:     "a/path",
			exitCode: 1, // dry run and verbose
		},
		{
			args:     []string{"prog", "-d", "-v", "--includeDotFiles", "a/path"},
			verbose:  true,
			dryRun:   true,
			dotFile:  true,
			path:     "a/path",
			exitCode: 1, // dry run and verbose
		},
	}

	var exitCode int
	exit = func(n int) {
		exitCode = n
	}

	for i, tt := range tests {
		exitCode = 0
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			os.Args = tt.args
			verbose, dryRun, dotFile, path := flagParse()
			if got, want := exitCode, tt.exitCode; got != want {
				t.Fatalf("exit got %d want %d", got, want)
			}
			if got, want := verbose, tt.verbose; got != want {
				t.Errorf("verbose got %t want %t", got, want)
			}
			if got, want := dryRun, tt.dryRun; got != want {
				t.Errorf("dryRun got %t want %t", got, want)
			}
			if got, want := dotFile, tt.dotFile; got != want {
				t.Errorf("path got %t want %t", got, want)
			}
			if got, want := path, tt.path; got != want {
				t.Errorf("path got %s want %s", got, want)
			}
		})
	}
}

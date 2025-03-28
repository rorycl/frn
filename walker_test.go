package main

import (
	"strings"
	"testing"
)

var expected string = strings.TrimSpace(`
[d] .
[d] A
[f]   A/%^&*()(___and
[f]   A/_AND
[d]   A/b
[f]     A/b/a nn $!@#
[d]     A/b/c d eFG
[d] b 1&2
[f]   b 1&2/12$-3.txt
[f]   b 1&2/12--3.txt
[f]   b 1&2/AnotherFile.Doc
`)

func TestWalker(t *testing.T) {
	tempDir := t.TempDir()
	// walker v1 : touch files
	err := walker("testdata", toucher(tempDir))
	if err != nil {
		t.Fatal(err)
	}
	b := strings.Builder{}
	testPrinter := printer(&b)

	// walker v2 : print files
	err = walker(tempDir, testPrinter)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(b.String()), expected; got != want {
		t.Errorf("got:\n%s\nwant\n%s\n", got, want)
	}
}

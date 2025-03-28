package main

import (
	"os"
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

func TestWalkerPrinter(t *testing.T) {
	b := strings.Builder{}
	testPrinter := printer(&b)
	err := walker("testdata", testPrinter)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(b.String()), expected; got != want {
		t.Errorf("got:\n%s\nwant\n%s\n", got, want)
	}
}

func TestWalkerToucher(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "walker_dir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatalf("cleanup err on directory removal: %v", err)
		}
	}()

	// walker v1 : touch files
	err = walker("testdata", toucher(tempDir))
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

/*
func main() {
	// err := walker("/home/rory/tmp", printer)
	err := walker("/home/rory/tmp/gcp_test2", toucher("/tmp/testdata2"))
	if err != nil {
		fmt.Println(err)
	}
}
*/

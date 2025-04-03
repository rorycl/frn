package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {

	// copy testdata to tempdir
	tempDir := t.TempDir()
	err := walker("testdata", toucher(tempDir))
	if err != nil {
		t.Fatal(err)
	}

	err = walker(tempDir, printer(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}

	// redirect output (normally os.Stdout)
	bb := &bytes.Buffer{}
	outputWriter = bb
	os.Args = []string{"prog", "-v", tempDir + "/"}

	main()
	fmt.Println(string(bb.Bytes()))

}

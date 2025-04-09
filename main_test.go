package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {

	// copy testdata to tempdir
	tempDir := t.TempDir()
	err := walker("testdata", toucher(tempDir))
	if err != nil {
		t.Fatal(err)
	}

	// expected
	want := `
          %^&*()(___and => and_and
          _AND => _and
            a nn $!@# => a_nn
          12$-3.txt => 12_3.txt
          12--n3.txt => 12_n3.txt
          AnotherFile.Doc => anotherfile.doc
            c d eFG => c_d_efg
          b => b
        A => a
        b 1&2 => b_1and2
`

	// redirect output (normally os.Stdout)
	bb := &bytes.Buffer{}
	outputWriter = bb
	os.Args = []string{"prog", "-v", tempDir + "/"}

	main()

	if got, want := string(bb.Bytes()), strings.TrimSpace(want); got != want {
		fmt.Errorf("got %s want %s", got, want)
	}

}

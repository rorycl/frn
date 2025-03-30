package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	tempDir := t.TempDir()
	// walker v1 : touch files
	err := walker("testdata", toucher(tempDir))
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile("", "frn_main_*")
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"prog", "-d", tempDir + "/"}
	out := os.Stdout
	os.Stdout = tempFile

	main()
	err = tempFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = out

	o, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	cont, err := ioutil.ReadAll(o)
	if err != nil {
		t.Fatal(err)
	}
	err = o.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("*****")
	fmt.Println(string(cont))
	os.Remove(tempFile.Name())

}

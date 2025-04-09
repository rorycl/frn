package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenamePath(t *testing.T) {

	fileRenamer = func(oldpath, newpath string) error {
		return nil
	}

	tests := []struct {
		origPath string
		isDir    bool
		newPath  string
		renamed  bool
		isErr    bool
	}{
		{
			origPath: "*(x[]Abc.doc",
			isDir:    false,
			newPath:  "x_abc.doc",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "ABCdef  g.doc",
			isDir:    false,
			newPath:  "abcdef_g.doc",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "$  #.doc",
			isDir:    false,
			newPath:  "_.doc",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "x& Y.doc",
			isDir:    false,
			newPath:  "xand_y.doc",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "abc",
			isDir:    true,
			newPath:  "abc",
			renamed:  false,
			isErr:    false,
		},
		{
			origPath: "abc", // check dirRegister
			isDir:    true,
			newPath:  "abc",
			renamed:  false,
			isErr:    false,
		},
		{
			origPath: "abc/deF",
			isDir:    true,
			newPath:  "abc/def",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "abc/", // empty
			isDir:    true,
			newPath:  "",
			renamed:  false,
			isErr:    false,
		},
		{
			origPath: "rename__test.go", // won't overwrite
			isDir:    false,
			newPath:  "rename_test.go",
			renamed:  true,
			isErr:    true,
		},
		{
			origPath: ".",
			isDir:    true,
			newPath:  ".",
			renamed:  false,
			isErr:    false,
		},
		{
			origPath: "nvim-linux-x86_64/share/nvim/runtime/lua/vim/func/_memoize.lua",
			isDir:    false,
			newPath:  "nvim-linux-x86_64/share/nvim/runtime/lua/vim/func/_memoize.lua",
			renamed:  false,
			isErr:    false,
		},
		{
			origPath: "/tmp/ABC_xyz.Doc", // capital in ext
			isDir:    false,
			newPath:  "/tmp/abc_xyz.doc",
			renamed:  true,
			isErr:    false,
		},
		{
			origPath: "/tmp/ abc_xyz.doc ", // spaces
			isDir:    false,
			newPath:  "/tmp/abc_xyz.doc",
			renamed:  true,
			isErr:    false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			newPath, renamed, err := pathRename(tt.origPath, tt.isDir)
			if got, want := newPath, tt.newPath; got != want {
				t.Errorf("path: got '%s' want '%s'", got, want)
			}
			if got, want := renamed, tt.renamed; got != want {
				t.Errorf("renamed: got %t want %t", got, want)
			}
			if got, want := err != nil, tt.isErr; got != want {
				t.Errorf("err: got %t want %t", got, want)
				fmt.Println(err)
			}
		})
	}
}

func TestRenameFunc(t *testing.T) {

	dir := t.TempDir()

	makeTestFile := func() string {
		o, err := os.Create(filepath.Join(dir, "_AND"))
		if err != nil {
			t.Fatal(err)
		}
		_ = o.Close()
		return o.Name()
	}

	tests := []struct {
		function renameFunc
		newFile  string // path
		output   string // for verbose/dry run output
		err      bool   //
	}{
		{
			function: printRename,
			newFile:  filepath.Join(dir, "a"),
			output:   fmt.Sprintf("%s => %s", "_AND", "a"),
			err:      false,
		},
		{
			function: wrappedOSRename,
			newFile:  filepath.Join(dir, "b"),
			output:   "",
			err:      false,
		},
		{
			function: verboseRename,
			newFile:  filepath.Join(dir, "c"),
			output:   fmt.Sprintf("%s => %s", "_AND", "c"),
			err:      false,
		},
	}

	/*
		listDir := func() {
			fs, err := ioutil.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range fs {
				fmt.Println(f.Name())
			}
		}
	*/

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			bb := &bytes.Buffer{}
			outputWriter = bb
			oldFile := makeTestFile()
			err := tt.function(oldFile, tt.newFile)
			if got, want := (err != nil), tt.err; got != want {
				t.Fatalf("unexpected error %v", err)
			}
			strResult := strings.TrimSpace(bb.String())
			if got, want := strResult, tt.output; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}

}

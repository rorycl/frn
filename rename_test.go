package main

import (
	"fmt"
	"testing"
)

func TestRename(t *testing.T) {

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
			origPath: "abc/", // empty
			isDir:    true,
			newPath:  "",
			renamed:  false,
			isErr:    false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			newPath, renamed, err := pathRename(tt.origPath, tt.isDir)
			if got, want := newPath, tt.newPath; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			if got, want := renamed, tt.renamed; got != want {
				t.Errorf("got %t want %t", got, want)
			}
			if got, want := err != nil, tt.isErr; got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}

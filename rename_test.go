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
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			newPath, renamed, err := pathRename(tt.origPath, tt.isDir)
			if got, want := newPath, tt.newPath; got != want {
				t.Errorf("path: got %s want %s", got, want)
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

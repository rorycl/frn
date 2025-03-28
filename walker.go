package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// printer is a closure returning a fs.WalkDirFunc for printing the
// contents of the path.
func printer(output io.Writer) fs.WalkDirFunc {
	output = output
	return func(path string, d fs.DirEntry, err error) error {
		indent := "  "
		countSlash := func(s string) int {
			return strings.Count(s, string(os.PathSeparator))
		}
		dirOrFile := func(d fs.DirEntry) string {
			if d == nil {
				return "x"
			}
			if d.IsDir() {
				return "d"
			}
			return "f"
		}
		_, err = fmt.Fprintf(output, "[%s] %s%s\n", dirOrFile(d), strings.Repeat(indent, countSlash(path)), path)
		return err
	}
}

// toucher is a closure returning an fs.WalDirFunc which recreates a
// file tree from path at target.
func toucher(target string) fs.WalkDirFunc {
	target = target
	return func(path string, d fs.DirEntry, err error) error {
		fullPath := filepath.Join(target, path)
		if d.IsDir() {
			if fullPath == target {
				return nil
			}
			return os.Mkdir(fullPath, 0755)
		}
		_, err = os.Create(fullPath)
		return err
	}
}

func walker(path string, fn fs.WalkDirFunc) error {
	dirFS := os.DirFS(path)
	return fs.WalkDir(dirFS, ".", fn)
}

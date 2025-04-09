package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func walker(path string, fn fs.WalkDirFunc) error {
	dirFS := os.DirFS(path)
	return fs.WalkDir(dirFS, ".", fn)
}

// toucher is a closure returning an fs.WalkDirFunc which recreates a
// file tree from path at target.
func toucher(target string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, _ error) error {
		fullPath := filepath.Join(target, path)
		if d.IsDir() {
			if fullPath == target {
				return nil
			}
			return os.Mkdir(fullPath, 0755)
		}
		_, err := os.Create(fullPath)
		return err
	}
}

// printer is a closure returning a fs.WalkDirFunc for printing the
// contents of the path.
func printer(output io.Writer, root string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, _ error) error {
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
		path = strings.ReplaceAll(path, root, "")
		_, err := fmt.Fprintf(output, "[%s] %s%s\n", dirOrFile(d), strings.Repeat(indent, countSlash(path)-1), path)
		return err
	}
}

var expected string = strings.TrimSpace(`
[f]   /A/%^&*()(___and
[f]   /A/_AND
[f]     /A/b/a nn $!@#
[f]   /b 1&2/12$-3.txt
[f]   /b 1&2/12--n3.txt
[f]   /b 1&2/AnotherFile.Doc
[d]     /A/b/c d eFG
[d]   /A/b
[d] /b 1&2
[d] /A

`)

func TestWalker(t *testing.T) {
	// setup copy
	tempDir := t.TempDir()
	err := walker("testdata", toucher(tempDir))
	if err != nil {
		t.Fatal(err)
	}

	// capture output
	b := strings.Builder{}
	testPrinter := printer(&b, tempDir)

	err = walkRename(tempDir, testPrinter)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(b.String()), expected; got != want {
		t.Errorf("got:\n%s\nwant\n%s\n", got, want)
	}
}

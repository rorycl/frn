package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

// toucher is a closure returning an fs.WalkDirFunc which recreates a
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

// walkRename walks the directory rooted at path, first applying
// pathRename for all files then for each directory sorted by longest
// path first.
func walkRename(path string) error {
	dirMap := map[string]int{}
	dirFS := os.DirFS(root)

	err := fs.WalkDir(dirFS, ".", func(path string, d fs.DirEntry, err error) error {
		p := filepath.Join(root, path)
		if p == root {
			return nil
		}
		if d.IsDir() {
			dirMap[p] = strings.Count(p, string(os.PathSeparator))
			return nil
		}
		_, _, err := pathRename(p, false)
		return err
	})
	if err != nil {
		return fmt.Errorf("file rename error: %v", err)
	}

	// sort directories by longest paths first
	dirs := make([]string, 0, len(dirMap))
	for d := range dirMap {
		dirs = append(dirs, d)
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirMap[dirs[i]] > dirMap[dirs[j]]
	})
	for _, p := range dirs {
		_, _, err := pathRename(p, false)
		if err != nil {
			return fmt.Errorf("directory rename error: %v", err)
		}
	}
	return nil
}

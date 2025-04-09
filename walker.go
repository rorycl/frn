package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// walkRename walks the directory rooted at path, first applying
// renameFunc for all files then for each directory sorted by longest
// path first.
func walkRename(path string, renameFunc fs.WalkDirFunc) error {
	type dirInfo struct {
		pathLen int         // number of os.PathSeparator
		d       fs.DirEntry // file/directory info
	}
	dirMap := map[string]dirInfo{}
	root := path
	dirFS := os.DirFS(root)

	err := fs.WalkDir(dirFS, ".", func(path string, d fs.DirEntry, _ error) error {
		p := filepath.Join(root, path)
		if p == root {
			return nil
		}
		if d.IsDir() {
			dirMap[p] = dirInfo{
				pathLen: strings.Count(p, string(os.PathSeparator)),
				d:       d,
			}
			return nil
		}
		return renameFunc(p, d, nil)
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
		if dirMap[dirs[i]].pathLen == dirMap[dirs[j]].pathLen {
			return dirMap[dirs[i]].d.Name() > dirMap[dirs[j]].d.Name()
		}
		return dirMap[dirs[i]].pathLen > dirMap[dirs[j]].pathLen
	})
	for _, p := range dirs {
		err := renameFunc(p, dirMap[p].d, nil)
		if err != nil {
			return fmt.Errorf("directory rename error: %v", err)
		}
	}
	return nil
}

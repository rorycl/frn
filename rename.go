package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ReplaceChars string = `[^A-Za-z0-9_.]`

var regexReplace = regexp.MustCompile(ReplaceChars)
var regexReplaceUnderscore = regexp.MustCompile("(_){2,}")

// dirRegister is a map of directory names as directories can be seen
// twice
var dirRegister = map[string]struct{}{}

var fileRenamer func(oldpath, newpath string) error = os.Rename

var noopRenamer func(string, string) error = func(oldPath, newPath string) error {
	indent := "  "
	countSep := func(s string) int {
		return strings.Count(s, string(os.PathSeparator))
	}
	fmt.Printf("%s%s => %s\n", strings.Repeat(indent, countSep(oldPath)), oldPath, newPath)
	return nil
}

// pathRename renames the file or directory at path returning the
// renamed filename, whether a rename occurred or error. The rename
// doesn't deal with odd characters in the extension.
//
// pathRename refuses to overwrite an existing file.
func pathRename(path string, isDir bool) (string, bool, error) {
	fileDir, fileName := filepath.Split(path)
	if fileName == "" {
		return "", false, nil
	}
	extension := filepath.Ext(fileName)
	nameSansExt := strings.TrimSuffix(fileName, extension)
	underFirstChar := strings.HasPrefix(nameSansExt, "_")

	newName := strings.ReplaceAll(nameSansExt, "&", "and")
	newName = regexReplace.ReplaceAllString(newName, "_")
	newName = strings.Trim(newName, "_")
	newName = strings.ToLower(newName)
	newName = regexReplaceUnderscore.ReplaceAllString(newName, "_")
	if newName == "" && extension != "." {
		newName = "_"
	}
	// put back leading underbar if it already existed
	if underFirstChar && !strings.HasPrefix(newName, "_") {
		newName = "_" + newName
	}

	newPath := filepath.Join(fileDir, newName) + extension

	if newPath == path {
		return newPath, false, nil
	}

	// directories can be seen twice
	if isDir {
		if _, ok := dirRegister[newPath]; ok {
			return newPath, false, nil
		}
		dirRegister[newPath] = struct{}{}
	}

	// don't overwrite.
	_, err := os.Stat(newPath)
	if err == nil {
		return newPath, true, fmt.Errorf("file %s already exists", newPath)
	}
	return newPath, true, fileRenamer(path, newPath)
}

// walkRename adapts a pathrename to a fs.WalkDirFunc
func walkRename(path string, d fs.DirEntry, _ error) error {
	_, _, err := pathRename(path, d.IsDir())
	return err
}

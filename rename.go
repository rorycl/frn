package main

import (
	"fmt"
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
	newName := strings.ReplaceAll(nameSansExt, "&", "and")
	newName = regexReplace.ReplaceAllString(newName, "_")
	newName = strings.Trim(newName, "_")
	newName = strings.ToLower(newName)
	newName = regexReplaceUnderscore.ReplaceAllString(newName, "_")
	if newName == "" {
		newName = "_"
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

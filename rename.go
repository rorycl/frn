package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ReplaceChars string = `[^A-Za-z0-9_.]`

var regexReplace = regexp.MustCompile(ReplaceChars)
var regexReplaceUnderscore = regexp.MustCompile("(_){2,}")

type renameFunc func(oldpath, newpath string) error

// fileRenamer is the func used to "rename" a file, but just potentially
// printing it or both renaming a file and printing it, etc.
//
// Any renameFunc provided here _must_ ensure that it gracefully deals
// with paths with an oldName the same as a newName without erroring.
var fileRenamer renameFunc

// wrappedOSRename is an os.Rename which returns nil if the old and new
// path are the same.
var wrappedOSRename renameFunc = func(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}
	return os.Rename(oldPath, newPath)
}

// default output is to os.Stdout
var outputWriter io.Writer = os.Stdout

// printRename only prints the old and new paths.
var printRename renameFunc = func(oldPath, newPath string) error {
	indent := "  "
	countSep := func(s string) int {
		return strings.Count(s, string(os.PathSeparator))
	}
	fmt.Fprintf(outputWriter, "%s%s => %s\n", strings.Repeat(indent, countSep(oldPath)), filepath.Base(oldPath), filepath.Base(newPath))
	return nil
}

// verboseRename does both a wrapped os.Rename and prints the change
var verboseRename renameFunc = func(oldPath, newPath string) error {
	err := wrappedOSRename(oldPath, newPath)
	if err != nil {
		return err
	}
	return printRename(oldPath, newPath)
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
	newName = strings.ToLower(newName)
	newName = regexReplaceUnderscore.ReplaceAllString(newName, "_")
	newName = strings.Trim(newName, "_")
	if newName == "" && extension != "." {
		newName = "_"
	}
	// put back leading underbar if it already existed
	if underFirstChar && !strings.HasPrefix(newName, "_") {
		newName = "_" + newName
	}

	ext := strings.ToLower(extension)
	ext = strings.TrimSpace(ext)

	newPath := filepath.Join(fileDir, newName) + ext

	renamed := (newPath != path)

	// don't overwrite.
	if renamed {
		_, err := os.Stat(newPath)
		switch {
		case err == nil && isDir:
			return "", false, nil
		case err == nil:
			return newPath, true, fmt.Errorf("file %s already exists", newPath)
		}
	}
	// fileRenamer _must_ handle not trying to rename a file or dir of
	// the same name
	return newPath, renamed, fileRenamer(path, newPath)
}

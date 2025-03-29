package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type processType int

const (
	NONE processType = iota
	FILE
	DIR
	WALK
)

// processKind (which isn't a great name) determines if the provided
// path exists, if it is a file, a directory, or a directory for
// transversal. These accord with the following path patterns:
//
//	// a file, type FILE
//	1. /here/filename
//	   ./here/filename
//	   filename
//
//	 // a directory only (or "naked" directory), type DIR
//	 2. /here/directory
//	    ./here/directory
//	    directory
//
//	 // a directory for walking, type WALK
//	 3. /here/directory/
//	    ./here/directory/
//	    directory/
//
// The difference between cases two and three is simply the final
// filepath.Separator in the third case.
func processKind(path string) (string, processType, error) {

	var pt processType
	hasTrailingSep := strings.HasSuffix(path, string(filepath.Separator))

	// note that os.Stat is _not_ run on the cleaned path since the path
	// provided to processKind is meaningful with a trailing Separator.
	info, err := os.Stat(path)
	if err != nil {
		return "", pt, fmt.Errorf("process stat error for path %s: %w", path, err)
	}
	isDir := info.IsDir()
	cleanedPath := filepath.Clean(path)

	switch {
	case !isDir:
		return cleanedPath, FILE, nil
	case isDir && !hasTrailingSep:
		return cleanedPath, DIR, nil
	default:
		return cleanedPath, WALK, nil
	}
}

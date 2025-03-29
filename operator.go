package main

// operator (which isn't a great name) determines if the provided path
// exists, if it is a file, a directory, or a directory for transversal.
// These accord with the following path patterns:
//
//	// a file, type FILE
//	1. /here/filename
//	   ./here/filename
//     filename
//
//  // a directory only (or "naked" directory), type DIR
//  2. /here/directory
//     ./here/directory
//     directory
//
//  // a directory for walking, type WALK
//  3. /here/directory/
//     ./here/directory/
//     directory/
//
// The difference between cases two and three is simply the final
// filepath.Separator in the third case.

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

func processKind(path string) (processType, error) {

	var pt processType
	hasTrailingSep := strings.HasSuffix(path, string(filepath.Separator))
	info, err := os.Stat(path)
	if err != nil {
		return pt, fmt.Errorf("process stat error for path %s: %w", path, err)
	}
	isDir := info.IsDir()

	switch {
	case !isDir:
		return FILE, nil
	case isDir && !hasTrailingSep:
		return DIR, nil
	default:
		return WALK, nil
	}
}

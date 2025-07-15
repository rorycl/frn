package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {

	// parse the command line flags.
	verbose, dryRun, incDotFiles, path := flagParse()

	// switch the fileRenamer func to either a print, os rename or
	// verbose os rename depending on the flags.
	switch {
	case dryRun:
		fileRenamer = printRename
	case verbose:
		fileRenamer = verboseRename
	default:
		fileRenamer = wrappedOSRename
	}

	checkErr := func(err error) {
		if err == nil {
			return
		}
		fmt.Println("error", err)
		os.Exit(1)
	}

	// determine what kind of processing is to be done.
	cleanPath, processType, err := processKind(path)
	checkErr(err)

	switch processType {
	case FILE:
		_, renamed, err := pathRename(cleanPath, false, incDotFiles)
		checkErr(err)
		if verbose && !renamed {
			fmt.Printf("%s didn't need renaming\n", path)
		}
	case DIR:
		_, renamed, err := pathRename(cleanPath, true, incDotFiles)
		checkErr(err)
		if verbose && !renamed {
			fmt.Printf("%s didn't need renaming\n", path)
		}
	case WALK: // recursive
		// walkPathRenameFunc adapts pathRename to a WalkDirFunc
		walkPathRenameFunc := func(path string, d fs.DirEntry, _ error) error {
			_, _, err := pathRename(path, d.IsDir(), incDotFiles)
			return err
		}
		err = walkRename(cleanPath, walkPathRenameFunc)
		checkErr(err)
	}
}

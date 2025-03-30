package main

import (
	"fmt"
	"os"
)

func main() {

	// parse the command line flags.
	verbose, dryRun, path := flagParse()
	if dryRun {
		verbose = true
		fileRenamer = noopRenamer // redirect file renaming func
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
		_, renamed, err := pathRename(cleanPath, false)
		checkErr(err)
		if verbose && !renamed {
			fmt.Printf("%s didn't need renaming\n", path)
		}
	case DIR:
		_, renamed, err := pathRename(cleanPath, true)
		checkErr(err)
		if verbose && !renamed {
			fmt.Printf("%s didn't need renaming\n", path)
		}
	case WALK: // recursive
		err = walker(cleanPath, walkRename)
		checkErr(err)
	}
}

// frn/filerenamer is a file renaming tool that works to remove
// unexpected characters from file and directory names, replacing
// sequences of the characters in ReplaceChars with the character "_".
// The program also lowercases file and directory names.
//
// As a special case the replacement sequence "_-_" is replaced with
// "_".
//
// frn works in three modes:
//
//  1. if provided with a single file, it will rename that file
//  2. if provided with a directory (without a trailing slash), it
//     will rename that directory
//  3. provided with a directory with a trailing slash it will rename the
//     directory and all contents below.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var ReplaceChars string = `\[\]-_:,&()+@ \`

var regexReplace = regexp.MustCompile("[" + ReplaceChars + "]+")

// Options are the command line options
type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"verbose: record changes"`
	DryRun  bool `short:"d" long:"dryrun" description:"dry-run mode: no changes will be made"`
	Args    struct {
		DirOrFilePath string `description:"directory path to process"`
	} `positional-args:"yes" required:"yes"`
}

type funcRenamer func(oldpath, newpath string) error

var FileRenamer funcRenamer = os.Rename

func printRenamer(oldpath, newpath string) error {
	fmt.Printf("rename %s -> %s\n", oldpath, newpath)
	return nil
}

var usage = fmt.Sprintf(`Path

Recursively rename the file, directory or the directory and all files
under it (by providing a directory ending with a "/") provided by Path.

The characters "%s" will be replaced by "_" and the names lowercased.
If in doubt run in test mode`, regexReplace)

var dirRegister = map[string]struct{}{}

// pathRenamer renames files and directories
func pathRenamer(path string, verbose, isDir bool) (fn string, err error) {

	fileDir, file := filepath.Split(path)

	if file == "" {
		if verbose {
			fmt.Printf("path %s has no file, returning early\n", fileDir)
		}
		return
	}

	newFile := regexReplace.ReplaceAllString(file, "_")
	newFile = strings.TrimRight(newFile, "_")
	newFile = strings.ToLower(newFile)
	newFile = strings.ReplaceAll(newFile, "_.", ".")
	fn = filepath.Join(fileDir, newFile)

	// directories can be seen twice
	if isDir {
		if _, ok := dirRegister[fn]; ok {
			if verbose {
				fmt.Printf("dir %s already seen, skipping...\n", fn)
			}
			return
		}
		dirRegister[fn] = struct{}{}
	}

	if verbose {
		fmt.Printf("dir %s old %s new %s skipping %t\n", fileDir, file, newFile, file == newFile)
	}
	if file == newFile {
		if verbose {
			fmt.Printf("not modifying %s, same name\n", fn)
		}
		return
	}

	_, err = os.Stat(fn)
	if err == nil {
		return fn, fmt.Errorf("file %s already exists", fn)
	}
	return fn, FileRenamer(path, fn)
}

func main() {

	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = usage

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// swap out to printer if in test mode
	if options.Test {
		FileRenamer = printRenamer
	}

	dirOrFile := filepath.Clean(options.Args.DirOrFilePath)

	info, err := os.Stat(dirOrFile)
	if err != nil {
		fmt.Printf("error reading input directory/file : %s\n", dirOrFile)
		os.Exit(1)
	}
	if options.Verbose {
		fmt.Printf("%s -> isdir() %v\n", dirOrFile, info.IsDir())
	}

	// rename root if necessary
	newDirOrFile, err := pathRenamer(dirOrFile, options.Verbose, info.IsDir())
	if err != nil {
		fmt.Printf("error renaming %s: %v\n", dirOrFile, err)
		os.Exit(1)
	}

	if !info.IsDir() {
		return
	}
	// process tree twice, first for files, second for directories
	for _, dirMode := range []bool{false, true} {
		filepath.WalkDir(newDirOrFile, func(path string, i fs.DirEntry, err error) error {
			if dirMode == i.IsDir() {
				_, err = pathRenamer(path, options.Verbose, i.IsDir())
				if err != nil {
					fmt.Printf("%s: error %s\n", path, err)
					os.Exit(1)
				}
				return err
			}
			return nil
		})
	}
}

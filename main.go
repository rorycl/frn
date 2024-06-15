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

var regexReplace = regexp.MustCompile("[-_,&()+@ ]+")

// Options are the command line options
type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"verbose"`
	Args    struct {
		DirOrFilePath string `description:"directory path to process"`
	} `positional-args:"yes" required:"yes"`
}

var usage = fmt.Sprintf(`DirPath

Recursively rename all files and directories at and under DirPath.
Renaming lowercases all files and directories and replaces the
characters "%s" with the replacement character "_".`, regexReplace)

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
	fn = fileDir + newFile

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
	return fn, os.Rename(fileDir+file, fn)
}

func main() {

	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = usage

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	dir := filepath.Clean(options.Args.DirPath)

	info, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("error reading dir: %s\n", dir)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("%s is not a directory", dir)
		os.Exit(1)
	}

	// rename root if necessary
	newDir, err := pathRenamer(dir, options.Verbose, info.IsDir())
	if err != nil {
		fmt.Printf("error renaming root %s: error %s\n", dir, err)
		os.Exit(1)
	}

	// process tree twice, first for files, second for directories
	for _, dirMode := range []bool{false, true} {
		filepath.WalkDir(newDir, func(path string, info fs.DirEntry, err error) error {
			if dirMode == info.IsDir() {
				_, err = pathRenamer(path, options.Verbose, info.IsDir())
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

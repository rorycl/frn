package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var usage string = `Path

Recursively rename the file, directory or the directory and all files
under it (by providing a directory ending with a "/") provided by Path.

All non-word, . and _ characters will be replaced by "_" and the names
lowercased. If in doubt run in dryrun mode.`

var exit func(int) = os.Exit

func flagParse() (verbose, dryRun bool, path string) {

	type options struct {
		Verbose bool `short:"v" long:"verbose" description:"verbose: record changes"`
		DryRun  bool `short:"d" long:"dryrun" description:"dry-run mode: no changes will be made"`
		Args    struct {
			DirOrFilePath string `description:"directory path to process"`
		} `positional-args:"yes" required:"yes"`
	}
	var opts options
	var parser = flags.NewParser(&opts, flags.Default)
	parser.Usage = usage

	if extraArgs, err := parser.Parse(); err != nil || len(extraArgs) > 0 {
		if len(extraArgs) > 0 {
			fmt.Printf("got unexpected additional arguments: %v\n", strings.Join(extraArgs, ","))
		}
		exit(1)
		return
	}
	if opts.Args.DirOrFilePath == "" {
		fmt.Println("no filepath found.")
		exit(1)
		return
	}
	if opts.DryRun && opts.Verbose {
		fmt.Println("dryrun and verbose selected -- please select one or ther other.")
		exit(1)

	}
	return opts.Verbose, opts.DryRun, opts.Args.DirOrFilePath
}

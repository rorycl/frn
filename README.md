# frn

**filerenamer**: rename files and directories, optionally recursively.

version 0.0.2 : 09 April 2025

Golang program for renaming files and directories, optionally
recursively, for \*nix operating systems.

The program works in three modes:

* file mode: only rename a file at the "end" of a path
* directory mode: only rename a directory -- use without a trailing
  slash
* recursive: rename everything under the target -- use with a trailing
  slash

All characters match the negative regexp `[^A-Za-z0-9_.]` are each
replaced with the `_` underbar character. Leading, trailing and
sequential underbar characters in the base directory or filename are
removed. Remaining characters, including those in filename extensions
are lower-cased. By default, dot files are ignored.

## Usage

```
Usage:
  frn Path

Recursively rename the file, directory or the directory and all files
under it (by providing a directory ending with a "/") provided by Path.

All non-word characters in the name of file (excluding the extension)
will be replaced by "_" and the names lowercased. 

Dotfiles are skipped unless the -i/--includeDotFiles option is included.

If in doubt run in dryrun mode. DirOrFilePath

Application Options:
  -v, --verbose          verbose: record changes
  -d, --dryrun           dry-run mode: no changes will be made
  -i, --includeDotFiles  also rename dot files

Help Options:
  -h, --help             Show this help message

Arguments:
  DirOrFilePath:         directory path to process

```

## Example

```
# example directory
> tree Abc
   Abc
   ├── AB d
   │   ├──  hi there.TXT
   │   └── There
   ├── c**D e
   │   └── eXample .JPG 
   └── fore.DOC

# rename only the directory, verbose mode
> ./frn -v Abc
   Abc => abc

# rename recursively
> ./frn abc/

# result
> tree abc
   abc
   ├── ab_d
   │   ├── hi_there.txt
   │   └── there
   ├── c_d_e
   │   └── example.jpg
   └── fore.doc

```


## License

Licensed under the [MIT Licence](LICENCE).

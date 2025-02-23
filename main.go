package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Dirsize processes a user provided, space-separated list of directories with optional flags and returns a list of their filesizes.
// TODO if I was working more on this:  unit tests, make some stylistic improvements, aim to reduce complexity, define more error types and cleaner error states

var (
	errorIncompleteArguments = fmt.Errorf("missing arguments")
)

func main() {
	if humanReadable, recursive, paths, err := processInput(os.Args); err == errorIncompleteArguments {
		printHelp()
		os.Exit(1)
	} else if sizeList, err := processPaths(humanReadable, recursive, paths); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	} else {
		if recursive {
			fmt.Println("(Recursive enabled)")
		}
		for path, size := range sizeList {
			fmt.Printf("%s: %s\n", path, size)
		}
	}
}

func processInput(input []string) (humanReadable, recursive bool, paths []string, err error) {
	// Ensure at least one argument has been provided, support a help flag
	if len(input) < 2 || input[1] == "--help" {
		return false, false, []string{}, errorIncompleteArguments
	}

	for _, arg := range input[1:] {
		if strings.Contains(arg, "--") {
			chunks := strings.Split(arg, "--")
			if chunks[1] == "human" {
				humanReadable = true
			} else if chunks[1] == "recursive" {
				recursive = true
			}
		} else {
			paths = append(paths, arg)
		}
	}
	return
}

func processPaths(humanReadable, recursive bool, paths []string) (sizeList map[string]string, err error) {
	sizeList = make(map[string]string, len(paths))
	for _, path := range paths {
		if absPath, err := filepath.Abs(path); err != nil {
			return map[string]string{}, err
		} else if size, err := getDirSize(recursive, absPath); err != nil {
			return map[string]string{}, err
		} else if humanReadable {
			sizeList[absPath] = formatSize(size)
		} else {
			sizeList[absPath] = strconv.FormatInt(size, 10)
		}
	}

	return sizeList, nil
}

// formatSize helper was generated by AI to save a little time since it's a straightforward math problem.
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	suffixes := []string{"KB", "MB", "GB"} // Note -- units can go higher, but I'm assuming GBs are a reasonable max for this.
	for n := size / unit; n >= unit && exp < len(suffixes)-1; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %s", float64(size)/float64(div), suffixes[exp])
}

// getDirSize does the heavy lifting for each path.
// If recursive is false, then it'll only process the top-level directory.  It's false by default unless the user toggles it on via input flag.
func getDirSize(recursive bool, path string) (int64, error) {
	var totalSize int64
	if info, err := os.Stat(path); err != nil {
		return 0, fmt.Errorf("failed to stat: %v", err)
	} else if !info.IsDir() { // Don't run on filenames, only directories
		return 0, fmt.Errorf("must provide a valid directory path")
	} else if err = filepath.WalkDir(path, func(entryPath string, d os.DirEntry, err error) error {
		return processEntry(entryPath, d, path, recursive, &totalSize)
	}); err != nil {
		return 0, fmt.Errorf("error walking directory: %v", err)
	} else {
		return totalSize, nil
	}
}

// processEntry processes each directory entry and updates totalSize accordingly
func processEntry(entryPath string, d os.DirEntry, rootPath string, recursive bool, totalSize *int64) error {
	if err := validateEntry(entryPath, d, rootPath, recursive); err != nil {
		return err
	} else if info, err := d.Info(); err != nil {
		return err
	} else if !info.IsDir() {
		*totalSize += info.Size()
	}

	return nil
}

// validateEntry checks if an entry should be processed or skipped
func validateEntry(entryPath string, d os.DirEntry, rootPath string, recursive bool) error {
	if entryPath == rootPath {
		return nil // Skip processing the root directory
	} else if d.IsDir() && !recursive {
		return filepath.SkipDir // Skip subdirectories if not a recursive request
	}
	return nil
}

func printHelp() {
	fmt.Println(`
--- Dirsize Help ---
Dirsize accepts a list of directory paths and returns their size.

Usage: dirsize [options] <dir1> <dir2> ...

Options:
  --help        Display this help content.
  --recursive   Include subdirectories in the total size.
  --human       Display sizes in human-readable format (KB, MB, GB).

Example:
  dirsize --human --recursive mydir anotherdir
`)
}

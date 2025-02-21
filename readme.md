To build the executable on Windows: go build -o dirsize.exe main.go
To build the executable on Mac/Linux: go build -o dirsize.sh main.go

Dirsize accepts a list of directory paths and returns their size.

Usage: dirsize [options] <dir1> <dir2> ...

Options:
--help        Display this help content.
--recursive   Include subdirectories in the total size.
--human       Display sizes in human-readable format (KB, MB, GB).

Example:
dirsize --human --recursive mydir anotherdir
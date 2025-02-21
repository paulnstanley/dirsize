
# Dirsize

This project is a utility which can read and output the size of a directory or directories.  It works with Mac/Linux and Windows.  




## Features

- Mac/Windows Support
- Optional recursion
- Optional human-readable sizes
- CLI-based




## Installation and Usage

Build on Windows:
```
go build -o dirsize.exe main.go
```

Build on Mac:
```
go build -o dirsize.sh main.go
```

Run on Windows:
```
dirsize [options] <dir1> <dir2> ...
```

Run on Mac:
```
./dirsize.sh [options] <dir1> <dir2> ...
```
    
## Documentation


```
Usage: dirsize [options] <dir1> <dir2> ...

--help        Display this help content.
--recursive   Include subdirectories in the total size.
--human       Display sizes in human-readable format (KB, MB, GB).

Example:
dirsize --human --recursive mydir anotherdir
```


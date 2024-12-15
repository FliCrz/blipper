# Repository Activity Score 
@author: kktkkk aka flicrz

## Description
This is the source code to create an activity score based on a CSV file.<br>
The application takes as possible inputs the CSV filepath, the number of days being analyzed, 
an optional filter to get a simple score based on it and eventually a "-d" debug flag.<br>
If no parameter is passed it will use default parameters and the Score Algorythm bellow.

## Score Algorythm
- Recently Updated => numberOfDays - scoreLastUpdate() (scoreLastUpdate creates an [][]int64 sorted descending by timestamp, 
where each slice represents a chunk of max and min timestamps and returns the index where the commit timestamp is in).
- Per File change => 10
- Per User => 5
- Per Commit => 2
- Per Addition or Deletion => 1

## Usage
After building (see bellow) simply open a terminal where the binary is located and run:
```
Usage:
    blipper
  or
    blipper [-f] [-n] [-t] [-d] [-v] [-h] 

  -f, --filename          CSV filename to read from (default: ../assets/commits.txt)
  -n, --numberOfDays      number of days being analyzed (default: 100, minimum 1)
  -t, --target            (optional) target to where to apply score if not set it will used our score algorythm
						  options: timestamp, files, additions, deletions, users, commits
  -d, --debug             logging (default: false)
  -v, --version           current version
  -h, --help              this help message
```

## Build
### Build requirements
go1.23.4

### Build
imply run in a terminal in "src" folder:
To build this application run: `go build -o ../bin/blipper main.go`
To build this application for a different architecture os system run: `GOARCH=<arch> GOOS=<os> go build -o ../bin/blipper .`
<br><br>
where arch can be:
- amd64
- arm64
- ...
<br><br>
and os can be:
- darwin
- linux
- windows
<br><br>
see more in golang docs: https://pkg.go.dev/go
<br>
Binary will be in "bin" folder.

## Testing 
Simply run in a terminal in "src" folder: `go test` <br>
NOTE: Unit tests have been written with help of AI Gemini - less work.
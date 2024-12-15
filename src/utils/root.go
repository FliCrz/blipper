package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/FliCrz/blipper/src/types"
)

func ErrorLogger(err error) {
	if err != nil {
		panic(err)
	}
}

func Debugger(m string, debug bool) {
	if debug {
		log.Println(m)
	}
}

func ParseInt(s string) (i int64) {
	i, err := strconv.ParseInt(s, 0, 64)
	ErrorLogger(err)
	return i
}

func help(f string, n int64) {
	fmt.Printf(`
	
Welcome to blipper

Usage:
    blipper
  or
    blipper [-f] [-n] [-t] [-d] [-v] [-h] 

  -f, --filename          CSV filename to read from (default: %s)
  -n, --numberOfDays      number of days being analyzed (default: %d, minimum 1)
  -t, --target            (optional) target to where to apply score if not set it will use our score algorythm
						  options: timestamp, files, additions, deletions, users, commits
  -d, --debug             logging (default: false)
  -v, --version           current version
  -h, --help              this help message

`, f, n)
}

func ParseArgs(filepath, scoringFilter, version string, numberOfDays int64, debug bool) (string, int64, string, bool) {
	for n := range os.Args {
		if n > 0 {
			Debugger(fmt.Sprintf("parsing args: %s", os.Args), debug)
			switch os.Args[n] {
			case "-h", "--help":
				help(filepath, numberOfDays)
				os.Exit(0)
			case "-f", "--filename":
				filepath = os.Args[n+1]
			case "-n", "--numberOfDays":
				numberOfDays = ParseInt(os.Args[n+1])
				if numberOfDays <= 0 {
					ErrorLogger(fmt.Errorf("number of days mus be bigger than 0"))
				}
			case "-t", "--target":
				scoringFilter = os.Args[n+1]
			case "-v", "--version":
				fmt.Println(version)
				os.Exit(0)
			case "-d", "--debug":
				debug = true
			}
		}
	}
	return filepath, numberOfDays, scoringFilter, debug
}

func ReadCsvToCommits(filepath string, debug bool) [][]string {
	Debugger(fmt.Sprintf("reading file: %s", filepath), debug)
	f, err := os.Open(filepath)
	ErrorLogger(err)
	defer f.Close()
	r := csv.NewReader(f)
	raw, err := r.ReadAll()
	ErrorLogger(err)
	return raw
}

func SortByScore(repos []types.Repository) {
	sort.SliceStable(repos, func(i int, j int) bool {
		return repos[i].Score > repos[j].Score
	})
}

func SortCommitsDecreasing(commits []types.Commit, filter string, debug bool) []types.Commit {
	Debugger(fmt.Sprintf("sorting commits decreasing by %s", filter), debug)
	ok := false
	for _, s := range []string{"timestamp", "files", "additions", "deletions", "users", "commits"} {
		if filter == s {
			ok = true
		}
	}
	if !ok {
		ErrorLogger(fmt.Errorf("TypeError: unsupported type: %s", filter))
	}
	sort.SliceStable(commits, func(i, j int) bool {
		e, err := commits[i].GetValue(filter)
		ErrorLogger(err)
		f, err := commits[j].GetValue(filter)
		ErrorLogger(err)
		return e > f
	})
	return commits
}

func GroupByRepository(c []types.Commit, debug bool) map[string]types.Repository {
	Debugger("group commits by repository", debug)
	repoMap := make(map[string]types.Repository)
	for _, i := range c {
		repo := types.Repository{
			Repository: i.Repository,
			Score:      1,
			Commits:    append(repoMap[i.Repository].Commits, i),
		}
		repoMap[i.Repository] = repo
	}
	return repoMap
}

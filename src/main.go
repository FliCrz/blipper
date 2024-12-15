package main

import (
	"fmt"

	"github.com/FliCrz/blipper/src/types"
	"github.com/FliCrz/blipper/src/utils"
)

// variables
var (
	filepath      = "../assets/commits.csv"
	scoringFilter string
	version             = "0.0.1"
	numberOfDays  int64 = 100
	repos         []types.Repository
	debug         = false
)

var (
	debugger = utils.Debugger
	parseInt = utils.ParseInt
)

func parseCommits(raw [][]string) (commits []types.Commit) {
	debugger("parsing commits", debug)
	for n, i := range raw {
		if n > 0 {
			user := i[1]
			if user == "" {
				user = "unknown"
			}
			commit := types.Commit{
				Timestamp:  parseInt(i[0]),
				User:       user,
				Repository: i[2],
				Files:      parseInt(i[3]),
				Additions:  parseInt(i[4]),
				Deletions:  parseInt(i[5]),
			}
			commits = append(commits, commit)
		}
	}
	return commits
}

func main() {
	debugger("STARTING", debug)
	filepath, numberOfDays, scoringFilter, debug = utils.ParseArgs(filepath, scoringFilter, version, numberOfDays, debug)
	msg := fmt.Sprintf("You are requesting scoring for file: %s for %d days", filepath, numberOfDays)

	raw := utils.ReadCsvToCommits(filepath, debug)
	commits := parseCommits(raw)

	if scoringFilter != "" {
		debugger(fmt.Sprintf("SCORING FILTER RECEIVED %s", scoringFilter), debug)
		repoMap := utils.GroupByRepository(commits, debug)
		msg = fmt.Sprintf("%s with filter: %s", msg, scoringFilter)
		fmt.Println(msg)
		for _, r := range repoMap {
			debugger(fmt.Sprintf("LOOPING REPOSITORY: %s", r.Repository), debug)
			r.Score = r.ScoreByFilter(scoringFilter)
			if !debug {
				r.Commits = nil
			}
			repos = append(repos, r)
		}
	} else {
		debugger("APPLYING SCORING ALGORYTHM", debug)
		fmt.Println(msg)
		sorted := utils.SortCommitsDecreasing(commits, "timestamp", debug)
		repoMap := utils.GroupByRepository(sorted, debug)
		for _, r := range repoMap {
			debugger(fmt.Sprintf("LOOPING REPOSITORY: %s", r.Repository), debug)
			for _, c := range r.Commits {
				r.Score = 100 - c.ScoreByLastUpdate(sorted[0].Timestamp, sorted[len(sorted)-1].Timestamp, numberOfDays)
			}
			for n, s := range []string{"files", "additions", "deletions", "users", "commits"} {
				debugger(fmt.Sprintf("LOOPING FILTER: %s", s), debug)
				if n == 0 {
					r.Score += r.ScoreByFilter(s) * 10
				} else if n == 3 {
					r.Score += r.ScoreByFilter(s) * 5
				} else if n == 4 {
					r.Score += r.ScoreByFilter(s) * 2
				} else {
					r.Score += r.ScoreByFilter(s) * 1
				}
			}
			if !debug {
				r.Commits = nil
			}
			repos = append(repos, r)
		}
	}

	debugger("SORTING BY SCORE", debug)
	utils.SortByScore(repos)
	fmt.Printf("\n%v\n", repos[0:9])
}

package types

import (
	"fmt"
)

// Repository ...
type Repository struct {
	Repository string   `json:"repository"`
	Score      int64    `json:"score"`
	Commits    []Commit `json:"commits"`
}

// Commit ...
type Commit struct {
	Timestamp  int64  `json:"timestamp"`
	User       string `json:"user"`
	Repository string `json:"repository"`
	Files      int64  `json:"files"`
	Additions  int64  `json:"additions"`
	Deletions  int64  `json:"deletions"`
}

func (r *Repository) ScoreByFilter(f string) int64 {
	s := 0
	c := r.Commits
	for _, i := range c {
		switch f {
		case "timestamp":
			s += int(i.Timestamp) // TODO improve this
		case "files":
			s += int(i.Files)
		case "additions":
			s += int(i.Additions)
		case "deletions":
			s += int(i.Deletions)
		case "users":
			s++
		case "commits":
			s = len(c)
		default:
			panic(fmt.Errorf("TypeError: unsupported type: %s", f))
		}
	}
	return int64(s)
}

func (c *Commit) GetValue(k string) (int64, error) {
	switch k {
	case "timestamp":
		return c.Timestamp, nil
	case "files":
		return c.Files, nil
	case "additions":
		return c.Additions, nil
	case "deletions":
		return c.Deletions, nil
	default:
		return 0, fmt.Errorf("TypeError: unsupported type: %s", k)
	}
}

func (c *Commit) ScoreByLastUpdate(last, first, numberOfDays int64) int64 {
	if numberOfDays <= 0 {
		numberOfDays = 100
	}
	duration := last - first
	chunkDuration := duration / numberOfDays
	for i := 0; i < int(numberOfDays); i++ {
		max := first + int64(i)*chunkDuration
		min := first + int64(i+1)*chunkDuration
		if max >= c.Timestamp && c.Timestamp >= min {
			return int64(i)
		}
	}
	return 0
}

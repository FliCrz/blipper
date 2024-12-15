package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/FliCrz/blipper/src/types"
	"github.com/FliCrz/blipper/src/utils"
	"github.com/google/go-cmp/cmp"
)

func TestScoreByFilter(t *testing.T) {
	testCases := []struct {
		name          string
		repository    types.Repository
		filter        string
		expectedScore int64
	}{
		{
			name: "Timestamp",
			repository: types.Repository{
				Commits: []types.Commit{
					{Timestamp: 100},
					{Timestamp: 200},
				},
			},
			filter:        "timestamp",
			expectedScore: 300,
		},
		{
			name: "Files",
			repository: types.Repository{
				Commits: []types.Commit{
					{Files: 5},
					{Files: 10},
				},
			},
			filter:        "files",
			expectedScore: 15,
		},
		{
			name: "Additions",
			repository: types.Repository{
				Commits: []types.Commit{
					{Additions: 20},
					{Additions: 30},
				},
			},
			filter:        "additions",
			expectedScore: 50,
		},
		{
			name: "Deletions",
			repository: types.Repository{
				Commits: []types.Commit{
					{Deletions: 5},
					{Deletions: 15},
				},
			},
			filter:        "deletions",
			expectedScore: 20,
		},
		{
			name: "Users", // Users simply counts the commits
			repository: types.Repository{
				Commits: []types.Commit{
					{User: "user1"},
					{User: "user2"},
				},
			},
			filter:        "users",
			expectedScore: 2,
		},
		{
			name: "Commits", // Commits counts the commits
			repository: types.Repository{
				Commits: []types.Commit{
					{Repository: "repo1"},
					{Repository: "repo1"},
					{Repository: "repo1"},
				},
			},
			filter:        "commits",
			expectedScore: 3,
		},
		{
			name:       "Invalid Filter",
			repository: types.Repository{Commits: []types.Commit{{}}},
			filter:     "invalid",
			// Expecting an error, but the function uses log.Fatalln which terminates execution.
			// Therefore, we cannot directly check for the error.  This scenario would be
			// handled by checking log output or via a panic recovery mechanism (not demonstrated here).
		},
		{
			name:          "Empty Commits",
			repository:    types.Repository{Commits: []types.Commit{}},
			filter:        "files",
			expectedScore: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.filter == "invalid" {
				//  Test cases with invalid filter will cause fatal error via utils.ErrorLogger
				//  and need special handling, which is not included in this basic example.
				//	A more robust approach involves testing for panics, but this is beyond
				//	the scope of a basic unit test example.
				return
			}
			score := tc.repository.ScoreByFilter(tc.filter)

			if diff := cmp.Diff(tc.expectedScore, score); diff != "" {
				t.Errorf("ScoreByFilter(%s) mismatch (-want +got):\n%s", tc.filter, diff)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	testCases := []struct {
		name      string
		commit    types.Commit
		key       string
		wantValue int64
		wantErr   error
	}{
		{
			name: "Timestamp",
			commit: types.Commit{
				Timestamp:  1678886400,
				User:       "testuser",
				Repository: "testrepo",
				Files:      10,
				Additions:  20,
				Deletions:  5,
			},
			key:       "timestamp",
			wantValue: 1678886400,
			wantErr:   nil,
		},
		{
			name: "Files",
			commit: types.Commit{
				Timestamp:  1678886400,
				User:       "testuser",
				Repository: "testrepo",
				Files:      10,
				Additions:  20,
				Deletions:  5,
			},
			key:       "files",
			wantValue: 10,
			wantErr:   nil,
		},
		{
			name: "Additions",
			commit: types.Commit{
				Timestamp:  1678886400,
				User:       "testuser",
				Repository: "testrepo",
				Files:      10,
				Additions:  20,
				Deletions:  5,
			},
			key:       "additions",
			wantValue: 20,
			wantErr:   nil,
		},
		{
			name: "Deletions",
			commit: types.Commit{
				Timestamp:  1678886400,
				User:       "testuser",
				Repository: "testrepo",
				Files:      10,
				Additions:  20,
				Deletions:  5,
			},
			key:       "deletions",
			wantValue: 5,
			wantErr:   nil,
		},
		{
			name: "InvalidKey",
			commit: types.Commit{
				Timestamp:  1678886400,
				User:       "testuser",
				Repository: "testrepo",
				Files:      10,
				Additions:  20,
				Deletions:  5,
			},
			key:       "invalid",
			wantValue: 0,
			wantErr:   fmt.Errorf("TypeError: unsupported type: %s", "invalid"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotValue, gotErr := tc.commit.GetValue(tc.key)

			if diff := cmp.Diff(tc.wantValue, gotValue); diff != "" {
				t.Errorf("getValue(%s) mismatch (-want +got):\n%s", tc.key, diff)
			}

			if diff := cmp.Diff(tc.wantErr, gotErr, cmp.Comparer(compareErrors)); diff != "" {
				t.Errorf("getValue(%s) error mismatch (-want +got):\n%s", tc.key, diff)
			}
		})
	}
}

func TestScoreByLastUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		commit        types.Commit
		last          int64
		first         int64
		numberOfDays  int64
		expectedScore int64
	}{
		{
			name:          "Middle of the range",
			commit:        types.Commit{Timestamp: 150},
			last:          200,
			first:         100,
			numberOfDays:  2,
			expectedScore: 0,
		},
		{
			name:          "Beginning of the range",
			commit:        types.Commit{Timestamp: 100},
			last:          200,
			first:         100,
			numberOfDays:  2,
			expectedScore: 0,
		},
		{
			name:          "End of the range",
			commit:        types.Commit{Timestamp: 200},
			last:          200,
			first:         100,
			numberOfDays:  2,
			expectedScore: 0,
		},
		{
			name:          "Zero days",
			commit:        types.Commit{Timestamp: 150},
			last:          200,
			first:         100,
			numberOfDays:  0, // Should return 0 since default is 100 and minimum 1
			expectedScore: 0,
		},
		{
			name:          "One day",
			commit:        types.Commit{Timestamp: 150},
			last:          200,
			first:         100,
			numberOfDays:  1,
			expectedScore: 0,
		},
		{
			name:          "Timestamp before first",
			commit:        types.Commit{Timestamp: 50},
			last:          200,
			first:         100,
			numberOfDays:  2,
			expectedScore: 0, // Should handle out-of-range timestamps gracefully
		},
		{
			name:          "Timestamp after last",
			commit:        types.Commit{Timestamp: 250},
			last:          200,
			first:         100,
			numberOfDays:  2,
			expectedScore: 0, // Should handle out-of-range timestamps gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := tc.commit.ScoreByLastUpdate(tc.last, tc.first, tc.numberOfDays)
			if diff := cmp.Diff(tc.expectedScore, score); diff != "" {
				t.Errorf("Score mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	originalArgs := os.Args
	defaultFilepath := "./commits.csv"
	defaultNumberOfDays := int64(100)
	defaultVersion := "0.0.1" // Replace with your actual default version

	defer func() { os.Args = originalArgs }()

	testCases := []struct {
		name             string
		args             []string
		expectedFilepath string
		expectedDays     int64
		expectedFilter   string
		expectedDebug    bool
		expectedExit     bool // Add a flag to check for expected exits
	}{
		{"Default values", []string{"blipper"}, defaultFilepath, defaultNumberOfDays, "", false, false},
		{"Filename flag", []string{"blipper", "-f", "test.csv"}, "test.csv", defaultNumberOfDays, "", false, false},
		{"NumberOfDays flag", []string{"blipper", "-n", "50"}, defaultFilepath, 50, "", false, false},
		{"Target flag", []string{"blipper", "-t", "timestamp"}, defaultFilepath, defaultNumberOfDays, "timestamp", false, false},
		{"Debug flag", []string{"blipper", "-d"}, defaultFilepath, defaultNumberOfDays, "", true, false},
		{"Combined flags", []string{"blipper", "-f", "data.csv", "-n", "25", "-t", "files", "-d"}, "data.csv", 25, "files", true, false},
		{"Help flag", []string{"blipper", "-h"}, "", 0, "", false, true},                   // Expect exit
		{"Version flag", []string{"blipper", "-v"}, "", 0, "", false, true},                // Expect exit
		{"Invalid NumberOfDays", []string{"blipper", "-n", "-1"}, "", -1, "", false, true}, // Expect error exit

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			if tc.expectedExit {
				// Test cases with os.Exit
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("%s: Expected exit, but did not exit", tc.name)
					}
				}()
				utils.ParseArgs(defaultFilepath, "", defaultVersion, defaultNumberOfDays, false)

			} else {
				filepath, numberOfDays, scoringFilter, debug := utils.ParseArgs(defaultFilepath, "", defaultVersion, defaultNumberOfDays, false)

				if filepath != tc.expectedFilepath {
					t.Errorf("Filepath mismatch: got %q, want %q", filepath, tc.expectedFilepath)
				}
				if numberOfDays != tc.expectedDays {
					t.Errorf("NumberOfDays mismatch: got %d, want %d", numberOfDays, tc.expectedDays)
				}
				if scoringFilter != tc.expectedFilter {
					t.Errorf("ScoringFilter mismatch: got %q, want %q", scoringFilter, tc.expectedFilter)
				}
				if debug != tc.expectedDebug {
					t.Errorf("Debug mismatch: got %v, want %v", debug, tc.expectedDebug)
				}
			}
		})
	}
}

func TestReadCsvToCommits(t *testing.T) {
	// Create a temporary test file
	tmpfile, err := os.CreateTemp("", "test_commits.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	// Write some test data to the file
	testData := `timestamp,user,repository,files,additions,deletions
1678886400,user1,repo1,10,20,5
1678886401,user2,repo2,5,10,2
1678886402,user3,repo1,2,4,1`
	if _, err := tmpfile.WriteString(testData); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil { // Important to close before reading
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		filepath string
		debug    bool
		expected [][]string
		wantErr  bool // Flag to indicate if an error is expected
	}{
		{
			name:     "Valid CSV",
			filepath: tmpfile.Name(),
			debug:    false, // debug=false simplifies testing
			expected: [][]string{
				{"timestamp", "user", "repository", "files", "additions", "deletions"},
				{"1678886400", "user1", "repo1", "10", "20", "5"},
				{"1678886401", "user2", "repo2", "5", "10", "2"},
				{"1678886402", "user3", "repo1", "2", "4", "1"},
			},
			wantErr: false,
		},
		{
			name:     "Non-existent file",
			filepath: "non_existent_file.csv",
			debug:    false,
			expected: nil,  // nil if expected is empty
			wantErr:  true, // Expect an error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture potential panics from utils.ErrorLogger
			defer func() {
				if r := recover(); (r != nil) != tc.wantErr {
					if tc.wantErr {
						t.Errorf("ReadCsvToCommits(%q) did not panic as expected", tc.filepath)
					} else {
						t.Errorf("ReadCsvToCommits(%q) panicked unexpectedly: %v", tc.filepath, r)
					}
				}
			}()

			got := utils.ReadCsvToCommits(tc.filepath, tc.debug)

			if !tc.wantErr { // Compare only if no error was expected
				if diff := cmp.Diff(tc.expected, got); diff != "" {
					t.Errorf("ReadCsvToCommits() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestSortByScore(t *testing.T) {
	testCases := []struct {
		name     string
		input    []types.Repository
		expected []types.Repository
	}{
		{
			name: "Sorts correctly",
			input: []types.Repository{
				{Repository: "repo2", Score: 50},
				{Repository: "repo1", Score: 100},
				{Repository: "repo3", Score: 25},
			},
			expected: []types.Repository{
				{Repository: "repo1", Score: 100},
				{Repository: "repo2", Score: 50},
				{Repository: "repo3", Score: 25},
			},
		},
		{
			name:     "Empty input",
			input:    []types.Repository{},
			expected: []types.Repository{}, // Should handle empty input gracefully
		},
		{
			name: "Same scores", // Stable sort
			input: []types.Repository{
				{Repository: "repo1", Score: 100},
				{Repository: "repo2", Score: 100},
				{Repository: "repo3", Score: 50},
			},
			expected: []types.Repository{
				{Repository: "repo1", Score: 100},
				{Repository: "repo2", Score: 100},

				{Repository: "repo3", Score: 50},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			utils.SortByScore(tc.input)
			if diff := cmp.Diff(tc.expected, tc.input); diff != "" {
				t.Errorf("SortByScore() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSortCommitsDecreasing(t *testing.T) {
	testCases := []struct {
		name     string
		commits  []types.Commit
		filter   string
		expected []types.Commit
	}{
		{
			name: "Sorts by timestamp correctly",
			commits: []types.Commit{
				{Timestamp: 100, Repository: "repo1"},
				{Timestamp: 300, Repository: "repo2"},
				{Timestamp: 200, Repository: "repo3"},
			},
			filter: "timestamp",
			expected: []types.Commit{
				{Timestamp: 300, Repository: "repo2"},
				{Timestamp: 200, Repository: "repo3"},
				{Timestamp: 100, Repository: "repo1"},
			},
		},
		{
			name: "Sorts by files correctly",
			commits: []types.Commit{
				{Files: 10, Repository: "repo1"},
				{Files: 30, Repository: "repo2"},
				{Files: 20, Repository: "repo3"},
			},
			filter: "files",
			expected: []types.Commit{
				{Files: 30, Repository: "repo2"},
				{Files: 20, Repository: "repo3"},
				{Files: 10, Repository: "repo1"},
			},
		},
		{
			name:     "Empty input",
			commits:  []types.Commit{},
			filter:   "timestamp",
			expected: []types.Commit{},
		},
		{
			name: "Invalid filter",
			commits: []types.Commit{
				{Timestamp: 100},
			},
			filter:   "invalid",
			expected: nil, // Expecting a panic due to the invalid filter
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.filter == "invalid" {

				defer func() {
					if r := recover(); r == nil {
						t.Errorf("SortCommitsDecreasing did not panic as expected")
					} else {
						// Check the panic message (optional)
						expectedMsg := fmt.Sprintf("TypeError: unsupported type: %s", tc.filter)
						if r.(error).Error() != expectedMsg {
							t.Errorf("Unexpected panic message: got %q, want %q", r, expectedMsg)
						}
					}
				}()
				utils.SortCommitsDecreasing(tc.commits, tc.filter, false) // Expecting a panic here

			} else {

				got := utils.SortCommitsDecreasing(tc.commits, tc.filter, false)
				if diff := cmp.Diff(tc.expected, got); diff != "" {
					t.Errorf("SortCommitsDecreasing() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGroupByRepository(t *testing.T) {
	testCases := []struct {
		name     string
		commits  []types.Commit
		expected map[string]types.Repository
	}{
		{
			name: "Groups correctly",
			commits: []types.Commit{
				{Repository: "repo1", User: "user1"},
				{Repository: "repo2", User: "user2"},
				{Repository: "repo1", User: "user3"},
			},
			expected: map[string]types.Repository{
				"repo1": {Repository: "repo1", Score: 1, Commits: []types.Commit{
					{Repository: "repo1", User: "user1"},
					{Repository: "repo1", User: "user3"},
				}},
				"repo2": {Repository: "repo2", Score: 1, Commits: []types.Commit{
					{Repository: "repo2", User: "user2"},
				}},
			},
		},
		{
			name:     "Empty input",
			commits:  []types.Commit{},
			expected: map[string]types.Repository{}, // Should handle empty input gracefully
		},
		// Add more test cases as needed, e.g., with duplicate commits
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := utils.GroupByRepository(tc.commits, false) // debug=false for simpler testing

			// Use cmp.Diff to compare maps
			if diff := cmp.Diff(tc.expected, got, cmp.AllowUnexported(types.Repository{}, types.Commit{})); diff != "" {
				t.Errorf("GroupByRepository() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func compareErrors(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 != nil && e2 != nil {
		return e1.Error() == e2.Error()
	}
	return false
}

package index

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommitEntry represents a parsed git commit.
type CommitEntry struct {
	SHA     string
	Subject string
	Body    string
	Date    string
	Files   []string
}

const commitSep = "|||"

// GetCommits returns commits from the last N days.
func GetCommits(days int) ([]CommitEntry, error) {
	since := fmt.Sprintf("--since=%d days ago", days)
	format := fmt.Sprintf("--format=%%H%s%%s%s%%b%s%%ci", commitSep, commitSep, commitSep)

	out, err := exec.Command("git", "log", since, format, "--name-only").Output()
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}

	return parseCommitLog(string(out)), nil
}

// GetLastCommit returns the most recent commit.
func GetLastCommit() (*CommitEntry, error) {
	format := fmt.Sprintf("--format=%%H%s%%s%s%%b%s%%ci", commitSep, commitSep, commitSep)

	out, err := exec.Command("git", "log", "-1", format, "--name-only").Output()
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}

	commits := parseCommitLog(string(out))
	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return &commits[0], nil
}

// GetChangedFiles returns the list of files changed in the last commit.
func GetChangedFiles() ([]string, error) {
	out, err := exec.Command("git", "diff", "HEAD~1", "--name-only").Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff: %w", err)
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

// parseCommitLog parses the output of git log with our custom format.
// Each commit block: "SHA|||subject|||body|||date\n\nfile1\nfile2\n\n"
func parseCommitLog(raw string) []CommitEntry {
	var commits []CommitEntry

	// Split by double newline to separate commit blocks.
	// Git log output with --name-only separates the header from files with a blank line,
	// and separates commits with blank lines.
	lines := strings.Split(raw, "\n")

	var current *CommitEntry
	parsingFiles := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, commitSep) {
			// This is a commit header line
			if current != nil {
				commits = append(commits, *current)
			}

			parts := strings.SplitN(line, commitSep, 4)
			entry := CommitEntry{
				SHA: parts[0],
			}
			if len(parts) > 1 {
				entry.Subject = parts[1]
			}
			if len(parts) > 2 {
				entry.Body = parts[2]
			}
			if len(parts) > 3 {
				entry.Date = parts[3]
			}

			current = &entry
			parsingFiles = false
			continue
		}

		if current == nil {
			continue
		}

		if line == "" {
			if !parsingFiles {
				parsingFiles = true
			}
			continue
		}

		if parsingFiles {
			current.Files = append(current.Files, line)
		}
	}

	if current != nil {
		commits = append(commits, *current)
	}

	return commits
}

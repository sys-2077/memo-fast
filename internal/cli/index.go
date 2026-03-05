package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/sys-2077/memo-fast/internal/config"
	"github.com/sys-2077/memo-fast/internal/index"
)

var (
	incremental bool
	dryRun      bool
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index project files and git commits",
	RunE:  runIndex,
}

func init() {
	indexCmd.Flags().BoolVar(&incremental, "incremental", false, "only index files changed in the last commit")
	indexCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print stats without sending to API")
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w\nRun 'memo-fast init' to create a config file.", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	collectionName := config.NormalizeCollectionName(filepath.Base(cwd))

	var files []index.FileEntry
	var commits []index.CommitEntry

	if incremental {
		// Incremental: only files changed in last commit
		changedPaths, err := index.GetChangedFiles()
		if err != nil {
			return fmt.Errorf("getting changed files: %w", err)
		}

		files, err = index.WalkPaths(changedPaths, cfg)
		if err != nil {
			return fmt.Errorf("reading changed files: %w", err)
		}

		lastCommit, err := index.GetLastCommit()
		if err != nil {
			if verbose {
				fmt.Printf("Warning: could not get last commit: %v\n", err)
			}
		} else {
			commits = []index.CommitEntry{*lastCommit}
		}
	} else {
		// Cold start: walk all project files
		files, err = index.Walk(cfg)
		if err != nil {
			return fmt.Errorf("walking project files: %w", err)
		}

		commits, err = index.GetCommits(cfg.Commits.WindowDays)
		if err != nil {
			if verbose {
				fmt.Printf("Warning: could not get commits: %v\n", err)
			}
		}
	}

	fmt.Printf("Found %d files, %d commits\n", len(files), len(commits))

	if dryRun {
		if verbose {
			for _, f := range files {
				fmt.Printf("  file: %s (%d bytes)\n", f.Path, len(f.Content))
			}
			for _, c := range commits {
				fmt.Printf("  commit: %s %s (%d files)\n", c.SHA[:minLen(len(c.SHA), 8)], c.Subject, len(c.Files))
			}
		}
		fmt.Println("Dry run complete. No data sent.")
		return nil
	}

	// Build payloads
	filePayloads := make([]index.FilePayload, len(files))
	for i, f := range files {
		filePayloads[i] = index.FilePayload{
			Path:    f.Path,
			Content: f.Content,
		}
	}

	commitPayloads := make([]index.CommitPayload, len(commits))
	for i, c := range commits {
		commitPayloads[i] = index.CommitPayload{
			SHA:     c.SHA,
			Subject: c.Subject,
			Body:    c.Body,
			Date:    c.Date,
			Files:   c.Files,
		}
	}

	req := index.IndexRequest{
		Files:   filePayloads,
		Commits: commitPayloads,
		Config: index.ConfigPayload{
			Collection: collectionName,
		},
	}

	fmt.Printf("Sending %d files, %d commits...\n", len(files), len(commits))
	start := time.Now()

	resp, err := index.Send(cfg.API.URL, cfg.API.Key, req)
	if err != nil {
		return fmt.Errorf("sending index request: %w", err)
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("Indexed %d files, %d commits (%d entities) in %.1fs\n",
		resp.IndexedFiles, resp.IndexedCommits, resp.Entities, elapsed)

	if err := reportServerErrors(resp.Errors); err != nil {
		return err
	}

	return nil
}

func reportServerErrors(errors []string) error {
	if len(errors) == 0 {
		return nil
	}
	fmt.Printf("[✗] Server reported %d indexing errors:\n", len(errors))
	for _, errMsg := range errors {
		fmt.Printf("  - %s\n", errMsg)
	}
	return fmt.Errorf("index completed with %d server-reported errors", len(errors))
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}

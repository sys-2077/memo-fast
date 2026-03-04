package index

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codexor/memo-fast/internal/config"
)

// FileEntry represents a discovered project file with its contents.
type FileEntry struct {
	Path    string
	Content string
}

// Walk discovers all project files matching the config criteria.
// It walks each target_dir, filters by valid_extensions, and skips
// ignore_patterns and hidden directories.
func Walk(cfg *config.Config) ([]FileEntry, error) {
	var entries []FileEntry

	for _, dir := range cfg.Index.TargetDirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue // skip non-existent target dirs
		}

		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // skip inaccessible entries
			}

			name := d.Name()

			// Skip hidden directories
			if d.IsDir() && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}

			// Skip ignored patterns (directories)
			if d.IsDir() {
				for _, pattern := range cfg.Index.IgnorePatterns {
					matched, _ := filepath.Match(pattern, name)
					if matched || name == pattern {
						return filepath.SkipDir
					}
				}
				return nil
			}

			// Skip ignored patterns (files)
			for _, pattern := range cfg.Index.IgnorePatterns {
				matched, _ := filepath.Match(pattern, name)
				if matched {
					return nil
				}
			}

			// Check valid extension
			ext := filepath.Ext(name)
			if !hasExtension(ext, cfg.Index.ValidExtensions) {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil // skip unreadable files
			}

			entries = append(entries, FileEntry{
				Path:    path,
				Content: string(content),
			})

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("walking directory %s: %w", dir, err)
		}
	}

	return entries, nil
}

// WalkPaths discovers files from a specific list of paths, filtering by config.
func WalkPaths(paths []string, cfg *config.Config) ([]FileEntry, error) {
	var entries []FileEntry

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}

		name := filepath.Base(path)

		// Check valid extension
		ext := filepath.Ext(name)
		if !hasExtension(ext, cfg.Index.ValidExtensions) {
			continue
		}

		// Skip ignored patterns
		skip := false
		for _, pattern := range cfg.Index.IgnorePatterns {
			matched, _ := filepath.Match(pattern, name)
			if matched {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		entries = append(entries, FileEntry{
			Path:    path,
			Content: string(content),
		})
	}

	return entries, nil
}

func hasExtension(ext string, valid []string) bool {
	for _, v := range valid {
		if ext == v {
			return true
		}
	}
	return false
}

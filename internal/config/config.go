package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure for memo-fast.
type Config struct {
	API        APIConfig     `yaml:"api"`
	Collection string        `yaml:"collection"`
	Index      IndexConfig   `yaml:"index"`
	Commits    CommitsConfig `yaml:"commits"`
}

// APIConfig holds the MCPize API connection settings.
type APIConfig struct {
	URL string `yaml:"url"`
	Key string `yaml:"key"`
}

// IndexConfig controls which files are discovered during indexing.
type IndexConfig struct {
	TargetDirs      []string `yaml:"target_dirs"`
	ValidExtensions []string `yaml:"valid_extensions"`
	IgnorePatterns  []string `yaml:"ignore_patterns"`
}

// CommitsConfig controls git history window.
type CommitsConfig struct {
	WindowDays int `yaml:"window_days"`
}

// DefaultConfigPath returns the conventional config file path relative to cwd.
func DefaultConfigPath() string {
	return filepath.Join(".mcp", "memo-fast", "config.yaml")
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return &cfg, nil
}

// Save writes the config as YAML to the given path, creating parent directories.
func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory %s: %w", dir, err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config %s: %w", path, err)
	}

	return nil
}

// NormalizeCollectionName converts a project directory name into a valid
// collection name: lowercase, replace - and . with _, prepend memo_.
func NormalizeCollectionName(name string) string {
	name = strings.ToLower(name)
	re := regexp.MustCompile(`[-.\s]+`)
	name = re.ReplaceAllString(name, "_")
	// Remove any character that is not alphanumeric or underscore
	clean := regexp.MustCompile(`[^a-z0-9_]`)
	name = clean.ReplaceAllString(name, "")
	if strings.HasPrefix(name, "memo_") {
		return name
	}
	return "memo_" + name
}

// DefaultIndex returns the default index configuration.
func DefaultIndex() IndexConfig {
	return IndexConfig{
		TargetDirs: []string{
			"src", "packages", "services", "config",
			"docs", "scripts", "lib", "internal", "cmd",
		},
		ValidExtensions: []string{
			".py", ".go", ".js", ".ts", ".yaml", ".yml",
			".json", ".md", ".sql", ".toml", ".rs", ".java",
		},
		IgnorePatterns: []string{
			"__pycache__", "node_modules", ".git", "*.pyc",
			".venv", "vendor", "dist", "build", ".next",
		},
	}
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultMCPizeConfigPath = "~/.mcpize/config.json"

type mcpizeConfig struct {
	Token string `json:"token"`
}

// ResolveAPIKey returns the MCPize token used by the CLI.
// Priority:
// 1) ~/.mcpize/config.json token
// 2) legacy key from project config (api.key)
func ResolveAPIKey(legacyKey string) (string, error) {
	if token, err := LoadMCPizeToken(); err == nil && strings.TrimSpace(token) != "" {
		return strings.TrimSpace(token), nil
	}

	legacyKey = strings.TrimSpace(legacyKey)
	if legacyKey != "" {
		return legacyKey, nil
	}

	return "", fmt.Errorf("MCPize API key not found. Run 'npx mcpize login' and verify ~/.mcpize/config.json")
}

// LoadMCPizeToken reads ~/.mcpize/config.json and returns token.
func LoadMCPizeToken() (string, error) {
	path := os.Getenv("MEMO_FAST_MCPIZE_CONFIG")
	if strings.TrimSpace(path) == "" {
		path = defaultMCPizeConfigPath
	}

	resolved, err := expandHome(path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("reading MCPize config %s: %w", resolved, err)
	}

	var cfg mcpizeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parsing MCPize config %s: %w", resolved, err)
	}

	token := strings.TrimSpace(cfg.Token)
	if token == "" {
		return "", fmt.Errorf("MCPize config %s has empty token", resolved)
	}

	return token, nil
}

func expandHome(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}

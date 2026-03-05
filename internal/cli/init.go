package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sys-2077/memo-fast/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive config wizard",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// 1. API Key (required - from MCPize dashboard)
	apiKey := promptRequired(reader, "memo-fast API key (from https://mcpize.com/settings)")

	// 2. Collection is derived from cwd basename (non-interactive)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	collection := config.NormalizeCollectionName(filepath.Base(cwd))

	cfg := &config.Config{
		API: config.APIConfig{
			URL: "https://api.memo-fast.dev",
			Key: apiKey,
		},
		Collection: collection,
		Index:      config.DefaultIndex(),
		Commits: config.CommitsConfig{
			WindowDays: 3,
		},
	}

	configPath := cfgFile
	if err := config.Save(configPath, cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	return nil
}

// promptRequired shows a question and keeps asking until a non-empty value is given.
func promptRequired(reader *bufio.Reader, question string) string {
	for {
		fmt.Printf("%s: ", question)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("  This field is required.")
	}
}

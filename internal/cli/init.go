package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sys-2077/memo-fast/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project config",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	collection := config.NormalizeCollectionName(filepath.Base(cwd))
	apiURL := os.Getenv("MEMO_FAST_API_URL")
	if apiURL == "" {
		apiURL = "https://api.memo-fast.dev"
	}

	cfg := &config.Config{
		API: config.APIConfig{
			URL: apiURL,
			Key: "", // legacy field; token is read from ~/.mcpize/config.json
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
	fmt.Println("MCPize token is read from ~/.mcpize/config.json (run: npx mcpize login)")
	return nil
}

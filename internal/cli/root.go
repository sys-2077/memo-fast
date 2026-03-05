package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sys-2077/memo-fast/internal/version"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "memo-fast",
	Short:   "Index your codebase for semantic memory",
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version.Version, version.Commit, version.BuildDate),
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".mcp/memo-fast/config.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

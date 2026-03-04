package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const hookMarker = "memo-fast"

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git post-commit hook",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git post-commit hook for auto-indexing",
	RunE:  runHookInstall,
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove git post-commit hook",
	RunE:  runHookUninstall,
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
	rootCmd.AddCommand(hookCmd)
}

func runHookInstall(cmd *cobra.Command, args []string) error {
	// Check .git/ exists
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (no .git directory found)")
	}

	// Find memo-fast binary path
	binPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding memo-fast binary path: %w", err)
	}

	// Resolve symlinks for a stable path
	binPath, err = filepath.EvalSymlinks(binPath)
	if err != nil {
		return fmt.Errorf("resolving binary path: %w", err)
	}

	hookPath := filepath.Join(".git", "hooks", "post-commit")

	hookContent := fmt.Sprintf(`#!/bin/sh
# memo-fast: auto-index on commit
%s index --incremental 2>/dev/null &
`, binPath)

	// If hook already exists, check if it already has memo-fast
	if data, err := os.ReadFile(hookPath); err == nil {
		if strings.Contains(string(data), hookMarker) {
			fmt.Println("Git hook already installed.")
			return nil
		}
		// Append to existing hook
		hookContent = string(data) + fmt.Sprintf("\n# memo-fast: auto-index on commit\n%s index --incremental 2>/dev/null &\n", binPath)
	}

	// Ensure hooks directory exists
	hooksDir := filepath.Join(".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("creating hooks directory: %w", err)
	}

	if err := os.WriteFile(hookPath, []byte(hookContent), 0o755); err != nil {
		return fmt.Errorf("writing hook file: %w", err)
	}

	fmt.Println("Git hook installed. Auto-indexing after each commit.")
	return nil
}

func runHookUninstall(cmd *cobra.Command, args []string) error {
	hookPath := filepath.Join(".git", "hooks", "post-commit")

	data, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No post-commit hook found.")
			return nil
		}
		return fmt.Errorf("reading hook file: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, hookMarker) {
		fmt.Println("No memo-fast hook found in post-commit.")
		return nil
	}

	// Remove memo-fast lines
	var kept []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, hookMarker) {
			continue
		}
		kept = append(kept, line)
	}

	remaining := strings.TrimSpace(strings.Join(kept, "\n"))

	// If only shebang or empty, remove the file entirely
	if remaining == "" || remaining == "#!/bin/sh" {
		if err := os.Remove(hookPath); err != nil {
			return fmt.Errorf("removing hook file: %w", err)
		}
	} else {
		if err := os.WriteFile(hookPath, []byte(remaining+"\n"), 0o755); err != nil {
			return fmt.Errorf("writing hook file: %w", err)
		}
	}

	fmt.Println("Git hook removed.")
	return nil
}

package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codexor/memo-fast/internal/config"
	"github.com/spf13/cobra"
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

	// 1. API URL
	apiURL := prompt(reader, "memo-fast API URL", "https://api.memo-fast.dev")

	// 2. API Key (optional, server may not require it)
	apiKey := prompt(reader, "API key (optional, press Enter to skip)", "")

	// 3. Vector backend
	backend := prompt(reader, "Vector backend (pinecone/qdrant)", "pinecone")

	// 4. Qdrant/Pinecone URL (optional, server may have it via env)
	qdrantURL := prompt(reader, "Vector DB URL (optional if server has credentials)", "")

	// 5. API Key for vector DB (optional)
	qdrantAPIKey := prompt(reader, "Vector DB API key (optional if server has credentials)", "")

	// 6. Collection name (derived from cwd basename)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	defaultCollection := config.NormalizeCollectionName(filepath.Base(cwd))
	collection := prompt(reader, "Collection name", defaultCollection)

	cfg := &config.Config{
		API: config.APIConfig{
			URL: apiURL,
			Key: apiKey,
		},
		Qdrant: config.QdrantConfig{
			URL:    qdrantURL,
			APIKey: qdrantAPIKey,
		},
		Backend:    backend,
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

// prompt shows a question with a default value and returns user input or default.
func prompt(reader *bufio.Reader, question, defaultVal string) string {
	fmt.Printf("%s [%s]: ", question, defaultVal)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
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

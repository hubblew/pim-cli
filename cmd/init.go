package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hubblew/pim/internal/agents"
	"github.com/hubblew/pim/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new PIM configuration",
	Long:  "Detect LLM tools, discover existing instruction files, and create a pim.yaml configuration.",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Detect LLM agents
	fmt.Println("Detecting CLI agents in your system...")
	tools := agents.DetectAgents()

	if len(tools) == 0 {
		fmt.Println("No CLI agents detected in your system.")
		fmt.Println("Currently supported tools: GitHub Copilot")
		return nil
	}

	fmt.Println("\nDetected agents:")
	for i, tool := range tools {
		fmt.Printf("%d. %s\n", i+1, tool)
	}

	// Step 2: Ask user to choose a tool
	var selectedTool string
	if len(tools) == 1 {
		selectedTool = tools[0]
		fmt.Printf("\nUsing detected tool: %s\n", selectedTool)
	} else {
		fmt.Print("\nSelect a tool (enter number): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)

		var idx int
		_, err = fmt.Sscanf(input, "%d", &idx)
		if err != nil || idx < 1 || idx > len(tools) {
			return fmt.Errorf("invalid selection")
		}
		selectedTool = tools[idx-1]
	}

	// Step 3: Ask for configuration file name
	fmt.Print("\nEnter configuration file name (default: pim.yaml): ")
	configName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	configName = strings.TrimSpace(configName)
	if configName == "" {
		configName = "pim.yaml"
	}

	// Step 4: Check if config file already exists
	if _, err := os.Stat(configName); err == nil {
		return fmt.Errorf("configuration file '%s' already exists", configName)
	}

	// Step 5: Ask for instructions folder name
	fmt.Print("\nEnter instructions folder name (default: ./instructions): ")
	instructionsDir, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	instructionsDir = strings.TrimSpace(instructionsDir)
	if instructionsDir == "" {
		instructionsDir = "./instructions"
	}

	// Step 6: Create instructions directory if it doesn't exist
	if err := os.MkdirAll(instructionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create instructions directory: %w", err)
	}
	fmt.Printf("Instructions directory: %s\n", instructionsDir)

	// Step 7: Look for existing instruction files
	existingFiles := discoverInstructionFiles(instructionsDir)

	// Step 8: Generate pim.yaml based on detected tool
	cfg, err := generateConfig(selectedTool, instructionsDir, existingFiles)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Step 9: Write config to file
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configName, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n✓ Configuration file '%s' created successfully!\n", configName)
	fmt.Printf("✓ Instructions directory: %s\n", instructionsDir)
	if len(existingFiles) > 0 {
		fmt.Printf("✓ Found %d existing instruction file(s)\n", len(existingFiles))
	}
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Review and edit %s\n", configName)
	fmt.Printf("  2. Add your instruction files to %s\n", instructionsDir)
	fmt.Printf("  3. Run 'pim install' to apply the configuration\n")

	return nil
}

// discoverInstructionFiles looks for existing instruction files
func discoverInstructionFiles(instructionsDir string) []string {
	var files []string

	// Common instruction file locations
	candidates := []string{
		filepath.Join(instructionsDir, "*.md"),
		"AGENTS.md",
		".github/copilot-instructions.md",
	}

	for _, pattern := range candidates {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		files = append(files, matches...)
	}

	// List files in instructions directory if it exists
	if entries, err := os.ReadDir(instructionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				filePath := filepath.Join(instructionsDir, entry.Name())
				// Avoid duplicates
				found := false
				for _, f := range files {
					if f == filePath {
						found = true
						break
					}
				}
				if !found {
					files = append(files, filePath)
				}
			}
		}
	}

	return files
}

// generateConfig creates a config based on the selected tool and discovered files
func generateConfig(tool, instructionsDir string, existingFiles []string) (*config.Config, error) {
	cfg := config.NewConfig()

	switch tool {
	case "GitHub Copilot":
		cfg.Targets = generateGitHubCopilotTarget(instructionsDir, existingFiles)
	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool)
	}

	return cfg, nil
}

// generateGitHubCopilotTarget generates target configuration for GitHub Copilot
func generateGitHubCopilotTarget(instructionsDir string, existingFiles []string) []config.Target {
	target := config.Target{
		Name:         "copilot-instructions",
		Output:       ".github/copilot-instructions.md",
		StrategyType: config.StrategyConcat,
		Include:      []string{},
	}

	// Add existing files to include list
	for _, file := range existingFiles {
		// Skip the output file itself
		if file == ".github/copilot-instructions.md" {
			continue
		}
		target.Include = append(target.Include, file)
	}

	// If no files were found, add example includes
	if len(target.Include) == 0 {
		target.Include = []string{
			filepath.Join(instructionsDir, "intro.md"),
			filepath.Join(instructionsDir, "coding-style.md"),
		}
	}

	return []config.Target{target}
}

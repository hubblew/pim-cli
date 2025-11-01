package cmd

import (
	"fmt"
	"os"

	"github.com/igor-vovk/lmpm/internal/config"
	"github.com/igor-vovk/lmpm/internal/installer"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages from sources to targets",
	Long:  `Fetch sources and copy specified files to target directories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := "lmpm.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = ".lmpm.yaml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found (lmpm.yaml or .lmpm.yaml)")
			}
		}

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		inst := installer.New(cfg)
		if err := inst.Install(); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pim",
	Short: "Language Model Package Manager",
	Long:  `A command-line utility for managing language models instruction files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no subcommand is specified
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pim.yaml)")
}

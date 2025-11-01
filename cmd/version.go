package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of LMPM",
	Long:  `Display the current version of LMPM.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LMPM v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

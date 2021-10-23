package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  `Print version and exit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "Version: %s\n", Version)
		fmt.Fprintf(os.Stdout, "BuiltOn: %s\n", Time)
	},
}

// Version is overridden at build time with semver
var Version = "<not-assigned>"

// Time is overridden at build time
var Time = "<not-assigned>"

func init() {
	rootCmd.AddCommand(versionCmd)
}

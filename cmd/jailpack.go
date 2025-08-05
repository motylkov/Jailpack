// Package cmd implements the main CLI entry point for Jailpack.
package cmd

import (
	"os"

	"jailpack/internal/commands"

	"github.com/spf13/cobra"
)

// Execute runs the root command for Jailpack.
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "jailpack",
		Short: "A tool for jail management",
		Long:  `Jailpack is a tool for FreeBSD jail management`,
	}

	rootCmd.AddCommand(commands.BuildCmd())
	rootCmd.AddCommand(commands.RunCmd())
	rootCmd.AddCommand(commands.ListCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

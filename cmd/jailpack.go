package cmd

import (
	"os"

	"jailpack/internal/commands"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "jailpack",
		Short: "A tool for jail management",
		Long:  `Jailpack is a tool for FreeBSD jail management and its philosophy: simplicity, stability, performance, security.`,
	}

	rootCmd.AddCommand(commands.BuildCmd())
	rootCmd.AddCommand(commands.RunCmd())
	rootCmd.AddCommand(commands.ListCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

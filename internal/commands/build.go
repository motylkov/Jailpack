// Package commands provides CLI commands for Jailpack.
package commands

import (
	"fmt"

	"jailpack/internal/build"

	"github.com/spf13/cobra"
)

// BuildCmd returns the build command for Jailpack.
func BuildCmd() *cobra.Command {
	var buildOutput string

	cmd := &cobra.Command{
		Use:   "build [path]",
		Short: "Create Cage from application",
		Long:  "Packages application into portable .cage.tar.gz",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			src := args[0]
			fmt.Printf("Building Cage from: %s\n", src)
			if err := build.CreateCage(src, buildOutput); err != nil {
				return fmt.Errorf("build error: %w", err)
			}
			fmt.Printf("Cage created: %s\n", buildOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&buildOutput, "output", "o", "app.cage.tar.gz", "Cage archive name")
	return cmd
}

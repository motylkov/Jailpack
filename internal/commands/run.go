// Package commands provides CLI commands for Jailpack.
package commands

import (
	"fmt"
	"path/filepath"

	"jailpack/internal/run"

	"github.com/spf13/cobra"
)

// RunCmd returns the run command for Jailpack.
func RunCmd() *cobra.Command {
	var runName, runIP string

	cmd := &cobra.Command{
		Use:   "run [cage]",
		Short: "Run Cage",
		Long:  "Extracts and runs jail from .cage.tar.gz",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cageFile := args[0]
			jailDir := filepath.Join("/usr/jails", runName)

			fmt.Printf("Extracting Cage: %s → %s\n", cageFile, jailDir)
			if err := run.ExtractAndRun(cageFile, jailDir, runName, runIP); err != nil {
				return fmt.Errorf("run error: %w", err)
			}

			fmt.Printf("Cage '%s' successfully started\n", runName)
			return nil
		},
	}

	cmd.Flags().StringVar(&runName, "name", "cage-app", "Jail name")
	cmd.Flags().StringVar(&runIP, "ip", "10.0.0.10", "Jail IP address")
	return cmd
}

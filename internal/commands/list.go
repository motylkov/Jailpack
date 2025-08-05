// Package commands provides CLI commands for Jailpack.
package commands

import (
	"fmt"

	"jailpack/internal/list"

	"github.com/spf13/cobra"
)

// ListCmd returns the list command for Jailpack.
func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List running jails",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("Running jails:")
			if err := list.ShowJails(); err != nil {
				fmt.Printf("Error listing jails: %v\n", err)
			}
		},
	}

	return cmd
}

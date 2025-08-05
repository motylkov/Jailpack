package commands

import (
	"fmt"

	"jailpack/internal/list"

	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List running jails",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📋 Running jails:")
			if err := list.ShowJails(); err != nil {
				fmt.Printf("Error listing jails: %v\n", err)
			}
		},
	}

	return cmd
}

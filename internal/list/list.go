// Package list provides functions for listing running jails.
package list

import (
	"fmt"
	"os/exec"
)

// ShowJails prints the list of running jails.
func ShowJails() error {
	out, err := exec.Command("jls").Output()
	if err != nil {
		return fmt.Errorf("jls command error: %w", err)
	}

	fmt.Print(string(out))
	return nil
}

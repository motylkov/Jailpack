// Package run provides functions for running Cage jails.
package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	dirPerm = 0750
	ipParts = 4
)

// ExtractAndRun extracts and runs jail from .cage.tar.gz.
func ExtractAndRun(cageFile, jailDir, runName, runIP string) error {
	// Validate inputs
	if !isValidPath(cageFile) {
		return fmt.Errorf("invalid cage file path: %s", cageFile)
	}
	if !isValidPath(jailDir) {
		return fmt.Errorf("invalid jail directory path: %s", jailDir)
	}
	if !isValidJailName(runName) {
		return fmt.Errorf("invalid jail name: %s", runName)
	}
	if !isValidIP(runIP) {
		return fmt.Errorf("invalid IP address: %s", runIP)
	}

	if err := os.MkdirAll(jailDir, dirPerm); err != nil {
		return fmt.Errorf("mkdirall: %w", err)
	}

	if err := exec.Command("tar", "-xzf", cageFile, "-C", jailDir).Run(); err != nil {
		return fmt.Errorf("extract tar: %w", err)
	}

	fmt.Printf("Starting jail: %s (IP: %s)\n", runName, runIP)
	jailArgs := []string{
		"-c", "name=" + runName,
		"path=" + jailDir,
		"ip4=" + runIP,
		"exec.start=/app-start.sh",
		"mount.devfs",
		"devfs_ruleset=4",
	}

	if err := exec.Command("jail", jailArgs...).Run(); err != nil { //nolint:gosec
		return fmt.Errorf("jail start: %w", err)
	}

	return nil
}

// isValidPath validates that the path is safe.
func isValidPath(path string) bool {
	return !strings.Contains(path, "..")
}

// isValidJailName validates jail name format.
func isValidJailName(name string) bool {
	if name == "" {
		return false
	}
	// Check for dangerous characters
	if strings.ContainsAny(name, ";&|`$()") {
		return false
	}
	return true
}

// isValidIP validates IP address format (basic check).
func isValidIP(ip string) bool {
	if ip == "" {
		return false
	}
	// Basic IP validation - check for dots and numbers
	parts := strings.Split(ip, ".")
	if len(parts) != ipParts {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
	}
	return true
}

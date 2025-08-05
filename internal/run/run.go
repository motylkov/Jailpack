package run

import (
	"fmt"
	"os"
	"os/exec"
)

// ExtractAndRun extracts and runs jail from .cage.tar.gz
func ExtractAndRun(cageFile, jailDir, runName, runIP string) error {
	if err := os.MkdirAll(jailDir, 0755); err != nil {
		return err
	}

	if err := exec.Command("tar", "-xzf", cageFile, "-C", jailDir).Run(); err != nil {
		return fmt.Errorf("extraction error: %w", err)
	}

	fmt.Printf("🚀 Starting jail: %s (IP: %s)\n", runName, runIP)
	jailArgs := []string{
		"-c", "name=" + runName,
		"path=" + jailDir,
		"ip4=" + runIP,
		"exec.start=/app-start.sh",
		"mount.devfs",
		"devfs_ruleset=4",
	}

	if err := exec.Command("jail", jailArgs...).Run(); err != nil {
		return fmt.Errorf("jail start error: %w", err)
	}

	return nil
}

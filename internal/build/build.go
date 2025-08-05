// Package build provides functions for building Cage archives.
package build

import (
	"archive/tar"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	randomStringLen = 6
	filePerm        = 0600
	dirPerm         = 0750
)

// CreateCage creates .cage.tar.gz from application directory.
func CreateCage(appDir, target string) error {
	// Validate input paths
	if !isValidPath(appDir) {
		return fmt.Errorf("invalid app directory path: %s", appDir)
	}
	if !isValidPath(target) {
		return fmt.Errorf("invalid target path: %s", target)
	}

	tmp := filepath.Join(os.TempDir(), "cage-build-"+randomString(randomStringLen))
	defer func() {
		if err := os.RemoveAll(tmp); err != nil {
			fmt.Fprintf(os.Stderr, "failed to remove temp dir: %v\n", err)
		}
	}()

	if err := createBaseFS(tmp); err != nil {
		return fmt.Errorf("createBaseFS: %w", err)
	}

	appDst := filepath.Join(tmp, "app")
	if err := copyDir(appDir, appDst); err != nil {
		return fmt.Errorf("copyDir: %w", err)
	}

	script := `#!/bin/sh
cd /app
echo "Starting application..."
# Auto-detect executable file
APP=$(ls *.py *.js *.go main 2>/dev/null | head -1)
if [ -n "$APP" ]; then
    case "$APP" in
        *.go) go run "$APP" ;;
        *.py) python "$APP" ;;
        *.js) node "$APP" ;;
        *)    ./"$APP" ;;
    esac
else
    echo "No executable application found"
    exit 1
fi
`
	if err := os.WriteFile(filepath.Join(tmp, "app-start.sh"), []byte(script), filePerm); err != nil {
		return fmt.Errorf("write app-start.sh: %w", err)
	}

	outFile, err := os.Create(target) //nolint:gosec
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	return archiveToCage(tmp, outFile)
}

func createBaseFS(root string) error {
	dirs := []string{"app", "bin", "sbin", "lib", "libexec", "usr/bin", "usr/sbin", "etc", "tmp", "dev"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), dirPerm); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}
	resolv := "nameserver 8.8.8.8\nnameserver 8.8.4.4"
	if err := os.WriteFile(filepath.Join(root, "etc", "resolv.conf"), []byte(resolv), filePerm); err != nil {
		return fmt.Errorf("write resolv.conf: %w", err)
	}
	return nil
}

func copyDir(src, dst string) error {
	if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir: %w", err)
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if err := os.MkdirAll(target, dirPerm); err != nil {
				return fmt.Errorf("mkdirall: %w", err)
			}
			return nil
		}
		data, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			return fmt.Errorf("readfile: %w", err)
		}
		if err := os.WriteFile(target, data, filePerm); err != nil {
			return fmt.Errorf("writefile: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}
	return nil
}

func archiveToCage(root string, writer *os.File) error {
	tw := tar.NewWriter(writer)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close tar writer: %v\n", err)
		}
	}()

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk: %w", err)
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("fileinfoheader: %w", err)
		}
		rel, _ := filepath.Rel(root, path)
		header.Name = rel
		if rel == "." {
			return nil
		}
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("writeheader: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			return fmt.Errorf("readfile: %w", err)
		}
		if _, err := tw.Write(data); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("walk: %w", err)
	}
	return nil
}

func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = alphanum[i%len(alphanum)]
	}
	return string(bytes)
}

// isValidPath validates that the path is safe and doesn't contain dangerous patterns.
func isValidPath(path string) bool {
	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return false
	}
	// Check for absolute paths that might be dangerous
	if filepath.IsAbs(path) && !strings.HasPrefix(path, "/tmp") && !strings.HasPrefix(path, "/usr/jails") {
		return false
	}
	return true
}

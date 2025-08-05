package build

import (
	"archive/tar"
	"io/fs"
	"os"
	"path/filepath"
)

// CreateCage creates .cage.tar.gz from application directory
func CreateCage(appDir, target string) error {
	tmp := filepath.Join(os.TempDir(), "cage-build-"+randomString(6))
	defer os.RemoveAll(tmp)

	if err := createBaseFS(tmp); err != nil {
		return err
	}

	appDst := filepath.Join(tmp, "app")
	if err := copyDir(appDir, appDst); err != nil {
		return err
	}

	script := `#!/bin/sh
cd /app
echo "🚀 Starting application..."
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
    echo "❌ No executable application found"
    exit 1
fi
`
	if err := os.WriteFile(filepath.Join(tmp, "app-start.sh"), []byte(script), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return archiveToCage(tmp, outFile)
}

// createBaseFS — minimal FS for jail
func createBaseFS(root string) error {
	dirs := []string{"app", "bin", "sbin", "lib", "libexec", "usr/bin", "usr/sbin", "etc", "tmp", "dev"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			return err
		}
	}
	resolv := "nameserver 8.8.8.8\nnameserver 8.8.4.4"
	return os.WriteFile(filepath.Join(root, "etc", "resolv.conf"), []byte(resolv), 0644)
}

// copyDir — recursive copying
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0755)
	})
}

// archiveToCage — archiving to .tar.gz
func archiveToCage(root string, writer *os.File) error {
	tw := tar.NewWriter(writer)
	defer tw.Close()

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		header.Name = rel
		if rel == "." {
			return nil
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tw.Write(data)
		return err
	})
}

// randomString — for temporary names
func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = alphanum[i%len(alphanum)]
	}
	return string(bytes)
}

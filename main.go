package main

import (
	"archive/tar"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	var (
		buildOutput string
		runName     string
		runIP       string
	)

	rootCmd := &cobra.Command{Use: "jailpack"}

	buildCmd := &cobra.Command{
		Use:   "build [path]",
		Short: "Создать Cage из приложения",
		Long:  "Упаковывает приложение в переносимый .cage.tar.gz",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			fmt.Printf("📦 Собираем Cage из: %s\n", src)
			if err := createCage(src, buildOutput); err != nil {
				return fmt.Errorf("ошибка сборки Cage: %w", err)
			}
			fmt.Printf("✅ Cage создан: %s\n", buildOutput)
			return nil
		},
	}
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "app.cage.tar.gz", "Имя Cage-архива")

	runCmd := &cobra.Command{
		Use:   "run [cage]",
		Short: "Запустить Cage",
		Long:  "Распаковывает и запускает jail из .cage.tar.gz",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cageFile := args[0]
			jailDir := filepath.Join("/usr/jails", runName)

			fmt.Printf("📦 Распаковка Cage: %s → %s\n", cageFile, jailDir)
			if err := os.MkdirAll(jailDir, 0755); err != nil {
				return err
			}
			if err := exec.Command("tar", "-xzf", cageFile, "-C", jailDir).Run(); err != nil {
				return fmt.Errorf("ошибка распаковки Cage: %w", err)
			}

			fmt.Printf("🚀 Запуск jail: %s (IP: %s)\n", runName, runIP)
			jailArgs := []string{
				"-c", "name=" + runName,
				"path=" + jailDir,
				"ip4=" + runIP,
				"exec.start=/app-start.sh",
				"mount.devfs",
				"devfs_ruleset=4",
			}
			if err := exec.Command("jail", jailArgs...).Run(); err != nil {
				return fmt.Errorf("ошибка запуска jail: %w", err)
			}

			fmt.Printf("🎉 Cage '%s' успешно запущен\n", runName)
			return nil
		},
	}
	runCmd.Flags().StringVar(&runName, "name", "cage-app", "Имя jail")
	runCmd.Flags().StringVar(&runIP, "ip", "10.0.0.10", "IP адрес jail")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Список запущенных jail",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📋 Запущенные jail:")
			out, _ := exec.Command("jls").Output()
			fmt.Print(string(out))
		},
	}

	rootCmd.AddCommand(buildCmd, runCmd, listCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// createCage создаёт .cage.tar.gz из директории приложения
func createCage(appDir, target string) error {
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
echo "🚀 Запуск приложения..."
# Автоопределение исполняемого файла
APP=$(ls *.py *.js *.go main 2>/dev/null | head -1)
if [ -n "$APP" ]; then
    case "$APP" in
        *.go) go run "$APP" ;;
        *.py) python "$APP" ;;
        *.js) node "$APP" ;;
        *)    ./"$APP" ;;
    esac
else
    echo "❌ Не найдено исполняемое приложение"
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

// createBaseFS — минимальная ФС для jail
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

// copyDir — рекурсивное копирование
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

// archiveToCage — архивирование в .tar.gz
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

// randomString — для временных имён
func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = alphanum[i%len(alphanum)]
	}
	return string(bytes)
}

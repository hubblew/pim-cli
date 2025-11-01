package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/igor-vovk/lmpm/internal/config"
)

type Installer struct {
	config *config.Config
}

func New(cfg *config.Config) *Installer {
	return &Installer{
		config: cfg,
	}
}

func (i *Installer) Install() error {
	tempDir, err := os.MkdirTemp("", "lmpm-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	sourceCache := make(map[string]string)

	for _, source := range i.config.Sources {
		sourceDir := filepath.Join(tempDir, source.Key)

		fmt.Printf("Fetching source '%s' from %s...\n", source.Key, source.URL)

		client := &getter.Client{
			Src:  source.URL,
			Dst:  sourceDir,
			Mode: getter.ClientModeDir,
		}

		if err := client.Get(); err != nil {
			return fmt.Errorf("failed to fetch source '%s': %w", source.Key, err)
		}

		sourceCache[source.Key] = sourceDir
	}

	for _, target := range i.config.Targets {
		fmt.Printf("Installing target '%s' to %s...\n", target.Name, target.Output)

		if err := os.MkdirAll(target.Output, 0755); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %w", target.Output, err)
		}

		for _, include := range target.Include {
			sourceDir, ok := sourceCache[include.Source]
			if !ok {
				return fmt.Errorf("source '%s' not found", include.Source)
			}

			for _, file := range include.Files {
				srcPath := filepath.Join(sourceDir, file)
				dstPath := filepath.Join(target.Output, file)

				if err := copyFile(srcPath, dstPath); err != nil {
					return fmt.Errorf("failed to copy file '%s': %w", file, err)
				}

				fmt.Printf("  ✓ %s\n", file)
			}
		}
	}

	fmt.Println("Installation complete!")
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

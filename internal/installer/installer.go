package installer

import (
	"fmt"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/hubble-works/pim/internal/config"
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
	tempDir, err := os.MkdirTemp("", "pim-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDirsByName := make(map[string]string)

	for _, source := range i.config.Sources {
		sourceDir := filepath.Join(tempDir, source.Name)

		fmt.Printf("Fetching source '%s' from %s...\n", source.Name, source.URL)

		client := &getter.Client{
			Src:  source.URL,
			Dst:  sourceDir,
			Mode: getter.ClientModeDir,
		}

		if err := client.Get(); err != nil {
			return fmt.Errorf("failed to fetch source '%s': %w", source.Name, err)
		}

		sourceDirsByName[source.Name] = sourceDir
	}

	for _, target := range i.config.Targets {
		fmt.Printf("Installing target '%s' to %s...\n", target.Name, target.Output)

		strategy := createStrategy(target.Strategy, target.Output)

		if err := strategy.Prepare(); err != nil {
			return err
		}
		defer strategy.Close()

		for _, include := range target.Include {
			sourceDir, ok := sourceDirsByName[include.Source]
			if !ok {
				return fmt.Errorf("source '%s' not found", include.Source)
			}

			for _, file := range include.Files {
				srcPath := filepath.Join(sourceDir, file)

				if err := strategy.AddFile(srcPath, file); err != nil {
					return fmt.Errorf("failed to add file '%s': %w", file, err)
				}

				fmt.Printf("  ✓ %s\n", file)
			}
		}
	}

	fmt.Println("Installation complete!")
	return nil
}

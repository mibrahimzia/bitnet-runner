package embedder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mibrahimzia/bitnet-runner/internal/utils"
	"github.com/mibrahimzia/bitnet-runner/native/bitnet-engine/windows"
)

// ExtractEngine unpacks the embedded engine files to the local runtime directory
func ExtractEngine() (string, error) {
	// 1. Get target directory (e.g., ~/.bitnet_runner/bin/windows)
	targetDir, err := utils.GetRuntimeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get runtime dir: %w", err)
	}

	// 2. Create directory if not exists
	if err := utils.EnsureDir(targetDir); err != nil {
		return "", fmt.Errorf("failed to create runtime dir: %w", err)
	}

	// 3. Get the embedded filesystem
	// Note: Since we are Windows-only for now, we call the windows package directly
	assets := windows.GetFS()

	// 4. Walk through embedded files and write them to disk
	err = fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Read file from memory
		data, err := assets.Open(path)
		if err != nil {
			return err
		}
		defer data.Close()

		// Create file on disk
		destPath := filepath.Join(targetDir, path)
		out, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}
		defer out.Close()

		// Copy data
		_, err = io.Copy(out, data)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("extraction failed: %w", err)
	}

	// Return the full path to the executable
	return filepath.Join(targetDir, "bitnet.exe"), nil
}
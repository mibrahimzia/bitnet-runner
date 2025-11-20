package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mibrahimzia/bitnet-runner/internal/utils"
)

// ScanModels looks for .gguf files in the models directory
func ScanModels() ([]ModelInfo, error) {
	// 1. Get the data directory
	appDir, err := utils.GetAppDataDir()
	if err != nil {
		return nil, err
	}

	modelsDir := filepath.Join(appDir, "models")
	
	// Ensure directory exists, create if not
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models dir: %w", err)
	}

	var models []ModelInfo

	// 2. Read directory
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, err
	}

	// 3. Iterate and filter
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".gguf") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(modelsDir, entry.Name())

		models = append(models, ModelInfo{
			ID:       entry.Name(),
			Name:     cleanName(entry.Name()),
			Filename: entry.Name(),
			FilePath: fullPath,
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}

	return models, nil
}

// cleanName removes extension and makes it readable
func cleanName(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return name
}
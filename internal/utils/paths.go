package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppName = "bitnet-runner"
)

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// GetAppDataDir returns the main application data directory
// Windows: C:\Users\Name\.bitnet_runner
func GetAppDataDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "."+AppName), nil
}

// GetRuntimeDir returns the directory where executable binaries will be extracted
func GetRuntimeDir() (string, error) {
	appDir, err := GetAppDataDir()
	if err != nil {
		return "", err
	}
	// We store the engine in a specific "bin" folder inside our app directory
	return filepath.Join(appDir, "bin", runtime.GOOS), nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
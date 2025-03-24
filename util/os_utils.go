package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetExecutableDir returns the directory containing the current executable
func GetExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsDarwin returns true if running on macOS
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// GetLibraryExtension returns the appropriate library extension for the current platform
func GetLibraryExtension() string {
	switch runtime.GOOS {
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default:
		return ".so"
	}
}

// NormalizePath normalizes a path according to the current platform
func NormalizePath(path string) string {
	if IsWindows() {
		return filepath.FromSlash(path)
	}
	return filepath.ToSlash(path)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir checks if a path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetTempDir returns the system's temporary directory
func GetTempDir() string {
	return os.TempDir()
}

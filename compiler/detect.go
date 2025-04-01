package compiler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	ErrInvalidCompilerPath = "invalid compiler path: %s"
	ErrNoCompilerFound     = "no supported compiler found"
	ErrNoWindowsCompiler   = "no supported compiler found on Windows"
	ErrUnsupportedOS       = "unsupported operating system: %s"
	ErrUnsupportedCompiler = "unsupported compiler type: %s"
	ErrCompilerNotFound    = "compiler not found: %s"
	ErrVersionCheckFailed  = "failed to get compiler version: %v"
)

// CompilerType represents the type of C++ compiler
type CompilerType string

const (
	CompilerGCC   CompilerType = "gcc"
	CompilerClang CompilerType = "clang"
	CompilerMSVC  CompilerType = "msvc"
	CompilerAuto  CompilerType = "auto"
)

// CompilerInfo contains information about the detected compiler
type CompilerInfo struct {
	Type         CompilerType
	Version      string
	Path         string
	IncludePaths []string
	EnvSetup     *CompilerEnvSetup
}

// CompilerEnvSetup contains information about how to set up the compiler's environment
type CompilerEnvSetup struct {
	SetupScript string   // Path to environment setup script (e.g., vcvarsall.bat for MSVC)
	SetupArgs   []string // Arguments for the setup script
	SetupCmd    string   // Command to run the setup script (e.g., "cmd /c" for MSVC)
}

// DetectCompiler determines the appropriate compiler based on the OS and user preference
func DetectCompiler(preferred CompilerType) (*CompilerInfo, error) {
	if preferred != CompilerAuto {
		return detectSpecificCompiler(preferred)
	}

	// Auto-detect based on OS
	switch runtime.GOOS {
	case "windows":
		return detectWindowsCompiler()
	case "linux", "darwin":
		return detectUnixCompiler()
	default:
		return nil, fmt.Errorf(ErrUnsupportedOS, runtime.GOOS)
	}
}

func detectSpecificCompiler(compiler CompilerType) (*CompilerInfo, error) {
	switch compiler {
	case CompilerGCC:
		return checkGCC()
	case CompilerClang:
		return checkClang()
	case CompilerMSVC:
		if runtime.GOOS != "windows" {
			return nil, fmt.Errorf("MSVC compiler is only supported on Windows")
		}
		return checkMSVC()
	default:
		return nil, fmt.Errorf(ErrUnsupportedCompiler, compiler)
	}
}

func detectWindowsCompiler() (*CompilerInfo, error) {
	// Try MSVC first
	if info, err := checkMSVC(); err == nil {
		return info, nil
	}

	// Try GCC/MinGW
	if info, err := checkGCC(); err == nil {
		return info, nil
	}

	return nil, errors.New(ErrNoWindowsCompiler)
}

func detectUnixCompiler() (*CompilerInfo, error) {
	// Try Clang first
	if info, err := checkClang(); err == nil {
		return info, nil
	}

	// Try GCC
	if info, err := checkGCC(); err == nil {
		return info, nil
	}

	return nil, errors.New(ErrNoCompilerFound)
}

func checkGCC() (*CompilerInfo, error) {
	// Try different possible GCC names based on OS
	compilerNames := []string{"g++", "gcc"}
	if runtime.GOOS == "windows" {
		compilerNames = append(compilerNames, "mingw32-g++", "x86_64-w64-mingw32-g++")
	}

	var path string
	var err error
	for _, name := range compilerNames {
		path, err = exec.LookPath(name)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf(ErrCompilerNotFound, "g++")
	}

	// Validate path is safe
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf(ErrInvalidCompilerPath, path)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, path, "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(ErrVersionCheckFailed, err)
	}

	return &CompilerInfo{
		Type:    CompilerGCC,
		Version: string(output),
		Path:    path,
	}, nil
}

func checkClang() (*CompilerInfo, error) {
	// Try different possible Clang names based on OS
	compilerNames := []string{"clang++", "clang"}
	if runtime.GOOS == "windows" {
		compilerNames = append(compilerNames, "llvm-clang++")
	}

	var path string
	var err error
	for _, name := range compilerNames {
		path, err = exec.LookPath(name)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf(ErrCompilerNotFound, "clang++")
	}

	// Validate path is safe
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf(ErrInvalidCompilerPath, path)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, path, "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(ErrVersionCheckFailed, err)
	}

	return &CompilerInfo{
		Type:    CompilerClang,
		Version: string(output),
		Path:    path,
	}, nil
}

func checkMSVC() (*CompilerInfo, error) {
	// First check if cl.exe is available
	path, err := exec.LookPath("cl.exe")
	if err != nil {
		return nil, fmt.Errorf(ErrCompilerNotFound, "cl.exe")
	}

	// Validate path is safe
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf(ErrInvalidCompilerPath, path)
	}

	// Get the version info from cl.exe
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, path)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(ErrVersionCheckFailed, err)
	}

	// Find include paths relative to cl.exe location
	includePaths := []string{}
	compilerDir := filepath.Dir(path)

	// Look for include directory in the same directory as cl.exe
	includeDir := filepath.Join(compilerDir, "include")
	if _, err := os.Stat(includeDir); err == nil {
		includePaths = append(includePaths, includeDir)
	}

	// Look for include directory in parent directory
	parentIncludeDir := filepath.Join(filepath.Dir(compilerDir), "include")
	if _, err := os.Stat(parentIncludeDir); err == nil {
		includePaths = append(includePaths, parentIncludeDir)
	}

	return &CompilerInfo{
		Type:         CompilerMSVC,
		Version:      string(output),
		Path:         path,
		IncludePaths: includePaths,
	}, nil
}

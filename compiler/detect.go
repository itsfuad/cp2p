package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func detectSpecificCompiler(compiler CompilerType) (*CompilerInfo, error) {
	switch compiler {
	case CompilerGCC:
		return checkGCC()
	case CompilerClang:
		return checkClang()
	case CompilerMSVC:
		return checkMSVC()
	default:
		return nil, fmt.Errorf("unsupported compiler type: %s", compiler)
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

	return nil, fmt.Errorf("no supported compiler found on Windows")
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

	return nil, fmt.Errorf("no supported compiler found")
}

func checkGCC() (*CompilerInfo, error) {
	cmd := exec.Command("g++", "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	path, err := exec.LookPath("g++")
	if err != nil {
		return nil, err
	}

	return &CompilerInfo{
		Type:    CompilerGCC,
		Version: string(output),
		Path:    path,
	}, nil
}

func checkClang() (*CompilerInfo, error) {
	cmd := exec.Command("clang++", "--version")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	path, err := exec.LookPath("clang++")
	if err != nil {
		return nil, err
	}

	return &CompilerInfo{
		Type:    CompilerClang,
		Version: string(output),
		Path:    path,
	}, nil
}

func findMSVCIncludePath(vsPath string) string {
	msvcPath := filepath.Join(vsPath, "VC\\Tools\\MSVC")
	entries, err := os.ReadDir(msvcPath)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(msvcPath, entry.Name(), "include")
		}
	}
	return ""
}

func findSDKIncludePath() string {
	sdkPath := "C:\\Program Files (x86)\\Windows Kits\\10\\Include"
	entries, err := os.ReadDir(sdkPath)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(sdkPath, entry.Name(), "ucrt")
		}
	}
	return ""
}

func checkMSVC() (*CompilerInfo, error) {
	// First check if cl.exe is available
	path, err := exec.LookPath("cl.exe")
	if err != nil {
		return nil, err
	}

	// Get the version info from cl.exe
	cmd := exec.Command("cl.exe")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Find Visual Studio path by looking for cl.exe's parent directory
	vsPath := ""
	dir := filepath.Dir(path)
	for {
		if dir == "" || dir == "." || dir == "/" {
			break
		}
		if strings.Contains(filepath.Base(dir), "Microsoft Visual Studio") {
			vsPath = dir
			break
		}
		dir = filepath.Dir(dir)
	}

	includePaths := []string{}
	var envSetup *CompilerEnvSetup
	if vsPath != "" {
		if msvcPath := findMSVCIncludePath(vsPath); msvcPath != "" {
			includePaths = append(includePaths, msvcPath)
		}
		if sdkPath := findSDKIncludePath(); sdkPath != "" {
			includePaths = append(includePaths, sdkPath)
		}

		// Set up MSVC environment configuration
		vcvarsall := filepath.Join(vsPath, "VC\\Auxiliary\\Build\\vcvarsall.bat")
		if _, err := os.Stat(vcvarsall); err == nil {
			envSetup = &CompilerEnvSetup{
				SetupScript: vcvarsall,
				SetupArgs:   []string{"x64"},
				SetupCmd:    "cmd /c",
			}
		}
	}

	return &CompilerInfo{
		Type:         CompilerMSVC,
		Version:      string(output),
		Path:         path,
		IncludePaths: includePaths,
		EnvSetup:     envSetup,
	}, nil
}

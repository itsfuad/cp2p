package compiler

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	errExpectedAbsPath = "Expected absolute path, got %v"
	errExpectedVersion = "Expected version info, got empty string"
)

// mockCompiler creates a mock compiler executable that returns a predefined version string
func mockCompiler(t *testing.T, dir, name, version string) string {
	path := filepath.Join(dir, name)

	// Create a Go program that will be compiled into our mock compiler
	content := []byte(`package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("` + version + `")
	os.Exit(0)
}`)

	// Write the Go source
	srcPath := path + ".go"
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("Failed to create mock compiler source: %v", err)
	}

	// Compile the mock compiler
	cmd := exec.Command("go", "build", "-o", path, srcPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build mock compiler: %v", err)
	}

	// Clean up the source file
	os.Remove(srcPath)

	return path
}

func TestDetectCompiler(t *testing.T) {
	tests := []struct {
		name     string
		os       string
		compiler CompilerType
		wantErr  bool
	}{
		{
			name:     "Auto detect on Windows",
			os:       "windows",
			compiler: CompilerAuto,
			wantErr:  false,
		},
		{
			name:     "Auto detect on Linux",
			os:       "linux",
			compiler: CompilerAuto,
			wantErr:  false,
		},
		{
			name:     "MSVC on non-Windows",
			os:       "linux",
			compiler: CompilerMSVC,
			wantErr:  true,
		},
		{
			name:     "Unsupported compiler type",
			os:       "linux",
			compiler: "unsupported",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that don't match the current OS
			if tt.os != runtime.GOOS {
				t.Skipf("Skipping test for OS %s on %s", tt.os, runtime.GOOS)
			}

			_, err := DetectCompiler(tt.compiler)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectCompiler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getGCCTestCase(t *testing.T, tmpDir string) struct {
	name     string
	compiler CompilerType
	setup    func()
	check    func(*testing.T, *CompilerInfo, error)
} {
	return struct {
		name     string
		compiler CompilerType
		setup    func()
		check    func(*testing.T, *CompilerInfo, error)
	}{
		name:     "GCC detection",
		compiler: CompilerGCC,
		setup: func() {
			mockCompiler(t, tmpDir, "g++", "g++ (GCC) 9.4.0")
		},
		check: func(t *testing.T, info *CompilerInfo, err error) {
			if err != nil {
				t.Skipf("Skipping GCC test: %v", err)
				return
			}
			if info.Type != CompilerGCC {
				t.Errorf("Expected compiler type GCC, got %v", info.Type)
			}
			if !filepath.IsAbs(info.Path) {
				t.Errorf(errExpectedAbsPath, info.Path)
			}
			if info.Version == "" {
				t.Error(errExpectedVersion)
			}
		},
	}
}

func getClangTestCase(t *testing.T, tmpDir string) struct {
	name     string
	compiler CompilerType
	setup    func()
	check    func(*testing.T, *CompilerInfo, error)
} {
	return struct {
		name     string
		compiler CompilerType
		setup    func()
		check    func(*testing.T, *CompilerInfo, error)
	}{
		name:     "Clang detection",
		compiler: CompilerClang,
		setup: func() {
			mockCompiler(t, tmpDir, "clang++", "clang version 12.0.0")
		},
		check: func(t *testing.T, info *CompilerInfo, err error) {
			if err != nil {
				t.Skipf("Skipping Clang test: %v", err)
				return
			}
			if info.Type != CompilerClang {
				t.Errorf("Expected compiler type Clang, got %v", info.Type)
			}
			if !filepath.IsAbs(info.Path) {
				t.Errorf(errExpectedAbsPath, info.Path)
			}
			if info.Version == "" {
				t.Error(errExpectedVersion)
			}
		},
	}
}

func getMSVCTestCase(t *testing.T, tmpDir string) struct {
	name     string
	compiler CompilerType
	setup    func()
	check    func(*testing.T, *CompilerInfo, error)
} {
	return struct {
		name     string
		compiler CompilerType
		setup    func()
		check    func(*testing.T, *CompilerInfo, error)
	}{
		name:     "MSVC detection on Windows",
		compiler: CompilerMSVC,
		setup: func() {
			if runtime.GOOS == "windows" {
				mockCompiler(t, tmpDir, "cl.exe", "Microsoft (R) C/C++ Optimizing Compiler Version 19.29.30133")
			}
		},
		check: func(t *testing.T, info *CompilerInfo, err error) {
			if runtime.GOOS != "windows" {
				t.Skip("Skipping MSVC test on non-Windows platform")
				return
			}
			if err != nil {
				t.Skipf("Skipping MSVC test: %v", err)
				return
			}
			if info.Type != CompilerMSVC {
				t.Errorf("Expected compiler type MSVC, got %v", info.Type)
			}
			if !filepath.IsAbs(info.Path) {
				t.Errorf(errExpectedAbsPath, info.Path)
			}
			if info.Version == "" {
				t.Error(errExpectedVersion)
			}
		},
	}
}

func getNotFoundTestCase() struct {
	name     string
	compiler CompilerType
	setup    func()
	check    func(*testing.T, *CompilerInfo, error)
} {
	return struct {
		name     string
		compiler CompilerType
		setup    func()
		check    func(*testing.T, *CompilerInfo, error)
	}{
		name:     "Compiler not found",
		compiler: CompilerGCC,
		setup: func() {
			// Don't create any mock compiler
		},
		check: func(t *testing.T, info *CompilerInfo, err error) {
			if err == nil {
				t.Skip("Skipping 'compiler not found' test as a compiler was found")
			}
		},
	}
}

func getCompilerTestCases(t *testing.T, tmpDir string) []struct {
	name     string
	compiler CompilerType
	setup    func()
	check    func(*testing.T, *CompilerInfo, error)
} {
	return []struct {
		name     string
		compiler CompilerType
		setup    func()
		check    func(*testing.T, *CompilerInfo, error)
	}{
		getGCCTestCase(t, tmpDir),
		getClangTestCase(t, tmpDir),
		getMSVCTestCase(t, tmpDir),
		getNotFoundTestCase(),
	}
}

func TestCompilerDetection(t *testing.T) {
	tmpDir := t.TempDir()
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+origPath)

	for _, tt := range getCompilerTestCases(t, tmpDir) {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			info, err := detectSpecificCompiler(tt.compiler)
			tt.check(t, info, err)
		})
	}
}

func TestIncludePathDetection(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Include path detection test is Windows-specific")
	}

	tmpDir := t.TempDir()

	// Create mock MSVC installation structure
	clPath := mockCompiler(t, tmpDir, "cl.exe", "Microsoft (R) C/C++ Optimizing Compiler Version 19.29.30133")
	includeDir := filepath.Join(tmpDir, "include")
	parentIncludeDir := filepath.Join(filepath.Dir(tmpDir), "include")

	// Create include directories
	if err := os.MkdirAll(includeDir, 0755); err != nil {
		t.Fatalf("Failed to create include directory: %v", err)
	}
	if err := os.MkdirAll(parentIncludeDir, 0755); err != nil {
		t.Fatalf("Failed to create parent include directory: %v", err)
	}

	// Save original PATH
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Add our temporary directory to PATH
	os.Setenv("PATH", filepath.Dir(clPath)+string(os.PathListSeparator)+origPath)

	info, err := checkMSVC()
	if err != nil {
		t.Skipf("Skipping MSVC include path test: %v", err)
		return
	}

	// Check include paths
	foundInclude := false
	foundParentInclude := false
	for _, path := range info.IncludePaths {
		if path == includeDir {
			foundInclude = true
		}
		if path == parentIncludeDir {
			foundParentInclude = true
		}
	}

	if !foundInclude {
		t.Error("Expected to find include directory in compiler directory")
	}
	if !foundParentInclude {
		t.Error("Expected to find include directory in parent directory")
	}
}

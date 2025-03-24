package compiler

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const fileName = "test.cpp"

func TestCompile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a simple C++ test file
	testFile := filepath.Join(tmpDir, fileName)
	testContent := `
#include <iostream>

extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with different compilers
	compilers := []CompilerType{CompilerGCC, CompilerClang, CompilerMSVC}
	for _, compilerType := range compilers {
		t.Run(string(compilerType), func(t *testing.T) {
			// Skip MSVC test on non-Windows platforms
			if compilerType == CompilerMSVC && os.Getenv("GOOS") != "windows" {
				t.Skip("MSVC tests only run on Windows")
			}

			// Detect compiler
			compiler, err := DetectCompiler(compilerType)
			if err != nil {
				t.Skipf("Compiler %s not available: %v", compilerType, err)
			}

			// Test compilation
			libPath, err := Compile(testFile, tmpDir, compiler)
			if err != nil {
				t.Fatalf("Compile() error = %v", err)
			}

			// Check if the library file exists
			if _, err := os.Stat(libPath); err != nil {
				t.Fatalf("Library file not created: %v", err)
			}
		})
	}
}

func TestCompileWithOptions(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a simple C++ test file without standard library includes
	testFile := filepath.Join(tmpDir, fileName)
	testContent := `
extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Detect compiler first to get include paths
	compiler, err := DetectCompiler(CompilerGCC)
	if err != nil {
		t.Fatalf("Failed to detect compiler: %v", err)
	}

	// Test with custom options
	opts := &CompileOptions{
		OptimizationLevel: "-O0",      // No optimization
		Debug:             true,       // Debug symbols
		IncludePaths:      []string{}, // Let the compiler handle include paths
		LibraryPaths:      []string{}, // Let the compiler handle library paths
	}

	// Test compilation with options
	libPath, err := CompileWithOptions(testFile, tmpDir, compiler, opts)
	if err != nil {
		t.Fatalf("CompileWithOptions() error = %v", err)
	}

	// Check if the library file exists
	if _, err := os.Stat(libPath); err != nil {
		t.Fatalf("Library file not created: %v", err)
	}
}

func TestDetectAvailableCompilers(t *testing.T) {
	compilers := []CompilerType{CompilerGCC, CompilerClang, CompilerMSVC}
	for _, compilerType := range compilers {
		t.Run(string(compilerType), testCompilerDetection(compilerType))
	}
	t.Run("Auto-detect", testAutoDetection)
}

func testCompilerDetection(compilerType CompilerType) func(*testing.T) {
	return func(t *testing.T) {
		if compilerType == CompilerMSVC && os.Getenv("GOOS") != "windows" {
			t.Logf("MSVC: Skipped (not available on %s)", os.Getenv("GOOS"))
			return
		}

		compiler, err := DetectCompiler(compilerType)
		if err != nil {
			t.Logf("%s: Not available (%v)", compilerType, err)
			return
		}

		logCompilerInfo(t, compiler)
	}
}

func testAutoDetection(t *testing.T) {
	compiler, err := DetectCompiler(CompilerAuto)
	if err != nil {
		t.Logf("Auto-detect: Failed (%v)", err)
		return
	}

	t.Logf("Auto-detect: Success")
	logCompilerInfo(t, compiler)
}

func logCompilerInfo(t *testing.T, compiler *CompilerInfo) {
	t.Logf("%s: Available", compiler.Type)
	t.Logf("  Path: %s", compiler.Path)
	t.Logf("  Version: %s", compiler.Version)
	if len(compiler.IncludePaths) > 0 {
		t.Logf("  Include paths:")
		for _, path := range compiler.IncludePaths {
			t.Logf("    - %s", path)
		}
	}
}

func TestBuildCompileCommand(t *testing.T) {
	tests := []struct {
		name     string
		compiler *CompilerInfo
		opts     *CompileOptions
		wantErr  bool
	}{
		{
			name: "GCC",
			compiler: &CompilerInfo{
				Type: CompilerGCC,
				Path: "g++",
			},
			opts:    DefaultCompileOptions(),
			wantErr: false,
		},
		{
			name: "Clang",
			compiler: &CompilerInfo{
				Type: CompilerClang,
				Path: "clang++",
			},
			opts:    DefaultCompileOptions(),
			wantErr: false,
		},
		{
			name: "MSVC",
			compiler: &CompilerInfo{
				Type: CompilerMSVC,
				Path: "cl.exe",
			},
			opts:    DefaultCompileOptions(),
			wantErr: false,
		},
		{
			name: "Invalid compiler",
			compiler: &CompilerInfo{
				Type: "invalid",
				Path: "invalid",
			},
			opts:    DefaultCompileOptions(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testBuildCompileCommandCase(tt))
	}
}

func testBuildCompileCommandCase(tt struct {
	name     string
	compiler *CompilerInfo
	opts     *CompileOptions
	wantErr  bool
}) func(*testing.T) {
	return func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, fileName)
		outputPath := filepath.Join(tmpDir, "test.dll")

		defer func() {
			if r := recover(); r != nil {
				if !tt.wantErr {
					t.Errorf("buildCompileCommand() panic = %v", r)
				}
			}
		}()

		args := buildCompileCommand(testFile, outputPath, tt.compiler, tt.opts)
		if tt.wantErr {
			t.Error("buildCompileCommand() should have panicked")
			return
		}

		if !slices.Contains(args, testFile) {
			t.Error("buildCompileCommand() missing source file")
		}
	}
}

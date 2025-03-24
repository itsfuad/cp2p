package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CompileOptions contains options for the compilation process
type CompileOptions struct {
	OptimizationLevel string
	Debug             bool
	IncludePaths      []string
	LibraryPaths      []string
}

// DefaultCompileOptions returns default compilation options
func DefaultCompileOptions() *CompileOptions {
	return &CompileOptions{
		OptimizationLevel: "-O2",
		Debug:             false,
		IncludePaths:      []string{},
		LibraryPaths:      []string{},
	}
}

// Compile compiles the C++ source file into a shared library
func Compile(sourceFile, outputDir string, compiler *CompilerInfo) (string, error) {
	opts := DefaultCompileOptions()
	opts.IncludePaths = compiler.IncludePaths
	return CompileWithOptions(sourceFile, outputDir, compiler, opts)
}

// CompileWithOptions compiles the C++ source file with custom options
func CompileWithOptions(sourceFile, outputDir string, compiler *CompilerInfo, opts *CompileOptions) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate output library name based on OS
	libName := generateLibraryName(sourceFile)
	outputPath := filepath.Join(outputDir, libName)

	// Build compilation command based on compiler type
	args := buildCompileCommand(sourceFile, outputPath, compiler, opts)

	// Execute compilation
	cmd := exec.Command(compiler.Path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("compilation failed: %v", err)
	}

	return outputPath, nil
}

func generateLibraryName(sourceFile string) string {
	baseName := filepath.Base(sourceFile)
	baseName = baseName[:len(baseName)-len(filepath.Ext(baseName))]

	switch runtime.GOOS {
	case "windows":
		return baseName + ".dll"
	case "darwin":
		return "lib" + baseName + ".dylib"
	default: // Linux and others
		return "lib" + baseName + ".so"
	}
}

func buildCompileCommand(sourceFile, outputPath string, compiler *CompilerInfo, opts *CompileOptions) []string {
	var args []string

	switch compiler.Type {
	case CompilerGCC:
		args = buildGCCCommand(sourceFile, outputPath, opts)
	case CompilerClang:
		args = buildClangCommand(sourceFile, outputPath, opts)
	case CompilerMSVC:
		args = buildMSVCCommand(sourceFile, outputPath, opts)
	default:
		panic(fmt.Sprintf("unsupported compiler type: %s", compiler.Type))
	}

	return args
}

func buildGCCCommand(sourceFile, outputPath string, opts *CompileOptions) []string {
	args := []string{
		"-shared",
		"-fPIC",
		opts.OptimizationLevel,
		"-o", outputPath,
	}

	if opts.Debug {
		args = append(args, "-g")
	}

	for _, include := range opts.IncludePaths {
		args = append(args, "-I"+include)
	}

	for _, lib := range opts.LibraryPaths {
		args = append(args, "-L"+lib)
	}

	args = append(args, sourceFile)
	return args
}

func buildClangCommand(sourceFile, outputPath string, opts *CompileOptions) []string {
	// Clang uses the same flags as GCC
	return buildGCCCommand(sourceFile, outputPath, opts)
}

func buildMSVCCommand(sourceFile, outputPath string, opts *CompileOptions) []string {
	args := []string{
		"/LD", // Create DLL
		"/O2", // Optimization level 2
		"/Fe:" + outputPath,
	}

	// Add include paths from compiler info
	for _, include := range opts.IncludePaths {
		args = append(args, "/I\""+include+"\"")
	}

	if opts.Debug {
		args = append(args, "/Zi")
	}

	for _, include := range opts.IncludePaths {
		args = append(args, "/I\""+include+"\"")
	}

	for _, lib := range opts.LibraryPaths {
		args = append(args, "/LIBPATH:"+lib)
	}

	args = append(args, sourceFile)
	return args
}

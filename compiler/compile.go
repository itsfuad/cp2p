package compiler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	// If compiler requires environment setup, create and run a setup script
	if compiler.EnvSetup != nil {
		// Create a batch file to set up the environment and run the compilation
		batchFile := filepath.Join(outputDir, "compile.bat")
		batchContent := fmt.Sprintf(`@echo off
call "%s" %s
"%s" %s
`, compiler.EnvSetup.SetupScript, strings.Join(compiler.EnvSetup.SetupArgs, " "),
			compiler.Path, strings.Join(args, " "))
		if err := os.WriteFile(batchFile, []byte(batchContent), 0644); err != nil {
			return "", fmt.Errorf("failed to create batch file: %v", err)
		}

		// Run the batch file
		// Validate paths are safe
		if !filepath.IsAbs(compiler.EnvSetup.SetupCmd) || !filepath.IsAbs(batchFile) {
			return "", fmt.Errorf("invalid command or batch file path")
		}

		ctx := context.Background()
		cmd := exec.CommandContext(ctx, compiler.EnvSetup.SetupCmd, batchFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("compilation failed: %v", err)
		}
		return outputPath, nil
	}

	// For compilers that don't need environment setup, run directly
	// Validate compiler path is safe
	if !filepath.IsAbs(compiler.Path) {
		return "", fmt.Errorf("invalid compiler path: %s", compiler.Path)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, compiler.Path, args...)
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
		"/MD", // Use multithreaded DLL runtime
		"/Fe:" + outputPath,
	}

	// Map optimization levels
	switch opts.OptimizationLevel {
	case "-O0":
		args = append(args, "/Od") // No optimization
	case "-O1":
		args = append(args, "/O1") // Minimize size
	case "-O2":
		args = append(args, "/O2") // Maximize speed
	case "-O3":
		args = append(args, "/O2") // MSVC doesn't have O3, use O2
	}

	if opts.Debug {
		args = append(args, "/Zi")
	}

	// Add include paths
	for _, include := range opts.IncludePaths {
		args = append(args, "/I\""+include+"\"")
	}

	// Add library paths
	for _, lib := range opts.LibraryPaths {
		args = append(args, "/LIBPATH:\""+lib+"\"")
	}

	args = append(args, sourceFile)
	return args
}

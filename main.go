package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cp2p/binding"
	"cp2p/compiler"
	"cp2p/config"
	"cp2p/parser"
	"cp2p/util"
)

var (
	inputFile   = flag.String("input", "", "Path to the C++ source file or project entry point")
	outputDir   = flag.String("output", "./bindings", "Output directory for generated bindings")
	compilerOpt = flag.String("compiler", "auto", "Compiler choice (gcc, clang, msvc, auto)")
	configFile  = flag.String("config", "", "Optional JSON config file (if not provided, will parse C++ file)")
)

func main() {
	flag.Parse()

	// Validate required flags
	if *inputFile == "" {
		fmt.Println("Error: --input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := util.NewLogger()

	// Detect compiler
	detectedCompiler, err := compiler.DetectCompiler(compiler.CompilerType(*compilerOpt))
	if err != nil {
		logger.Fatalf("Failed to detect compiler: %v", err)
	}

	// Parse config or C++ file
	var cfg *config.Config
	if *configFile != "" {
		cfg, err = config.ParseConfig(*configFile)
		if err != nil {
			logger.Fatalf("Failed to parse config file: %v", err)
		}
	} else {
		cfg, err = parser.ParseCppFile(*inputFile)
		if err != nil {
			logger.Fatalf("Failed to parse C++ file: %v", err)
		}
	}

	// Compile C++ code
	libPath, err := compiler.Compile(*inputFile, *outputDir, detectedCompiler)
	if err != nil {
		logger.Fatalf("Failed to compile C++ code: %v", err)
	}

	// Generate Python bindings
	moduleName := filepath.Base(*inputFile)
	moduleName = moduleName[:len(moduleName)-len(filepath.Ext(moduleName))]

	if err := binding.GenerateBindings(moduleName, libPath, *outputDir, cfg); err != nil {
		logger.Fatalf("Failed to generate Python bindings: %v", err)
	}

	logger.Info(fmt.Sprintf("Successfully generated Python bindings in %s", *outputDir))
}

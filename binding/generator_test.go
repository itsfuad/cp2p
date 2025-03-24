package binding

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cp2p/config"
)

func TestGenerateBindings(t *testing.T) {
	// Create a temporary directory for test output
	tmpDir := t.TempDir()

	// Create a test config
	testConfig := &config.Config{
		Functions: []config.FunctionConfig{
			{
				Name:        "add",
				Description: "Adds two integers",
				Parameters: []config.Param{
					{Name: "a", Type: "int", Description: "First integer"},
					{Name: "b", Type: "int", Description: "Second integer"},
				},
				ReturnType: "int",
			},
			{
				Name:        "multiply",
				Description: "Multiplies two floating-point numbers",
				Parameters: []config.Param{
					{Name: "a", Type: "double", Description: "First number"},
					{Name: "b", Type: "double", Description: "Second number"},
				},
				ReturnType: "double",
			},
		},
	}

	// Test generating bindings
	err := GenerateBindings("test", "test.dll", tmpDir, testConfig)
	if err != nil {
		t.Fatalf("GenerateBindings() error = %v", err)
	}

	// Check if the output file exists
	outputPath := filepath.Join(tmpDir, "test.py")
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("Output file not created: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Check for expected content
	expectedStrings := []string{
		"import ctypes",
		"import sys",
		"import os",
		"from typing import Any, Union, Optional, List, Dict, Tuple",
		"TYPE_MAPPING = {",
		"'int': ctypes.c_int",
		"'double': ctypes.c_double",
		"_lib = None",
		"def add(a: int, b: int) -> int:",
		"def multiply(a: float, b: float) -> float:",
		"__all__ = ['add', 'multiply']",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(string(content), expected) {
			t.Errorf("Generated file missing expected content: %s", expected)
		}
	}
}

func TestNewGenerator(t *testing.T) {
	// Create a temporary directory for test output
	tmpDir := t.TempDir()

	// Create a test config
	testConfig := &config.Config{
		Functions: []config.FunctionConfig{
			{
				Name:        "add",
				Description: "Adds two integers",
				Parameters: []config.Param{
					{Name: "a", Type: "int", Description: "First integer"},
					{Name: "b", Type: "int", Description: "Second integer"},
				},
				ReturnType: "int",
			},
		},
	}

	// Test generating bindings using NewGenerator
	err := GenerateBindings("test", "test.dll", tmpDir, testConfig)
	if err != nil {
		t.Fatalf("GenerateBindings() error = %v", err)
	}

	// Check if the output file exists
	outputPath := filepath.Join(tmpDir, "test.py")
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
}

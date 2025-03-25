package binding

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"

	"cp2p/config"
)

// Generator handles the generation of Python bindings
type Generator struct {
	moduleName string
	libPath    string
	outputDir  string
	config     *config.Config
}

// NewGenerator creates a new binding generator
func NewGenerator(moduleName, libPath, outputDir string, cfg *config.Config) *Generator {
	return &Generator{
		moduleName: moduleName,
		libPath:    libPath,
		outputDir:  outputDir,
		config:     cfg,
	}
}

// GenerateBindings generates Python bindings for the C++ library
func GenerateBindings(moduleName, libPath, outputDir string, cfg *config.Config) error {
	gen := NewGenerator(moduleName, filepath.Base(libPath), outputDir, cfg)
	return gen.generate()
}

func (g *Generator) generate() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate the Python binding file
	outputPath := filepath.Join(g.outputDir, g.moduleName+".py")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Generate the binding code
	if err := g.generateBindingCode(file); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateBindingCode(file *os.File) error {
	// Define the template for the Python binding using html/template for security
	tmpl := template.Must(template.New("binding").Parse(pythonBindingTemplate))

	// Define type mappings
	typeMappings := map[string]string{
		"int":         "ctypes.c_int",
		"float":       "ctypes.c_float",
		"double":      "ctypes.c_double",
		"char":        "ctypes.c_char",
		"bool":        "ctypes.c_bool",
		"void":        "None",
		"const char*": "ctypes.c_char_p",
	}

	pythonTypeHints := map[string]string{
		"int":         "int",
		"float":       "float",
		"double":      "float",
		"char":        "str",
		"bool":        "bool",
		"void":        "None",
		"const char*": "str",
	}

	// Prepare template data
	data := struct {
		ModuleName      string
		LibPath         string
		Functions       []config.FunctionConfig
		Platform        string
		Types           []config.TypeConfig
		TypeMappings    map[string]string
		PythonTypeHints map[string]string
	}{
		ModuleName:      g.moduleName,
		LibPath:         g.libPath,
		Functions:       g.config.Functions,
		Platform:        runtime.GOOS,
		Types:           g.config.Types,
		TypeMappings:    typeMappings,
		PythonTypeHints: pythonTypeHints,
	}

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to generate binding code: %v", err)
	}

	return nil
}

// pythonBindingTemplate is the template for generating Python bindings
const pythonBindingTemplate = `import ctypes
import sys
import os
from typing import Any, Union, Optional, List, Dict, Tuple

# Basic type mapping (always included)
TYPE_MAPPING = {
    {{range $key, $value := .TypeMappings}}
    '{{$key}}': {{$value}},
    {{end}}
}

# Python type hints mapping (always included)
PYTHON_TYPE_HINTS = {
    {{range $key, $value := .PythonTypeHints}}
    '{{$key}}': '{{$value}}',
    {{end}}
}

{{range .Types}}
{{if eq .Kind "struct"}}
class {{.Name}}(ctypes.Structure):
    """
    {{.Description}}
    """
    _fields_ = [
        {{range .Fields}}
        ("{{.Name}}", TYPE_MAPPING["{{.Type}}"]),  # {{.Description}}
        {{end}}
    ]
{{else if eq .Kind "enum"}}
class {{.Name}}(ctypes.c_int):
    """
    {{.Description}}
    """
    {{range .Values}}
    {{.}} = {{.}}
    {{end}}
{{else if eq .Kind "union"}}
class {{.Name}}(ctypes.Union):
    """
    {{.Description}}
    """
    _fields_ = [
        {{range .Fields}}
        ("{{.Name}}", TYPE_MAPPING["{{.Type}}"]),  # {{.Description}}
        {{end}}
    ]
{{end}}

{{end}}

# Load the shared library based on the OS
_lib = None
if sys.platform.startswith('win'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), '{{.LibPath}}'))
elif sys.platform.startswith('linux'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), '{{.LibPath}}'))
elif sys.platform.startswith('darwin'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), '{{.LibPath}}'))

{{range .Functions}}
# Configure function signature for {{.Name}}
_lib.{{.Name}}.argtypes = [{{range $i, $p := .Parameters}}{{if $i}}, {{end}}TYPE_MAPPING["{{$p.Type}}"]{{end}}]
_lib.{{.Name}}.restype = TYPE_MAPPING["{{.ReturnType}}"]

def {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.Name}}: {{index $.PythonTypeHints $p.Type}}{{end}}) -> {{index $.PythonTypeHints .ReturnType}}:
    """
    {{.Description}}
    {{if .Docstring}}
    {{.Docstring}}
    {{end}}
    {{range .Parameters}}
    Args:
        {{.Name}} ({{index $.PythonTypeHints .Type}}): {{.Description}}
    {{end}}
    Returns:
        {{index $.PythonTypeHints .ReturnType}}: {{.Description}}
    """
    return _lib.{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.Name}}{{end}})

{{end}}

__all__ = [{{range $i, $f := .Functions}}{{if $i}}, {{end}}'{{$f.Name}}'{{end}}]
`

# C++ to Python Bindings Generator

A tool to generate Python bindings for C++ code, supporting multiple compilers and platforms.

## Features

- Automatic compiler detection (GCC, Clang, MSVC)
- Cross-platform support (Windows, Linux, macOS)
- Generates Python bindings using pybind11
- Handles C++ class and function bindings
- Supports custom include paths and compiler flags
- Configurable through JSON or C++ file annotations

## How It Works

1. **Compiler Detection**:
   - Automatically detects available C++ compilers in the system PATH
   - Supports GCC, Clang, and MSVC (Windows only)
   - Verifies compiler version and capabilities

2. **C++ Code Analysis**:
   - Parses C++ source files to identify classes and functions
   - Supports both automatic detection and manual configuration
   - Handles inheritance and virtual functions

3. **Binding Generation**:
   - Creates pybind11-based Python bindings
   - Generates proper type conversions
   - Handles memory management and object lifetimes
   - Supports both static and dynamic library generation

4. **Build Process**:
   - Compiles C++ code into a shared library
   - Links against Python and pybind11
   - Handles platform-specific build requirements

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/cp2p.git
cd cp2p

# Install dependencies
go mod download

# Build the tool
go build -o cp2p
```

## Usage

### Basic Usage

```bash
# Generate bindings from a C++ file
cp2p --input example.cpp --output ./bindings

# Use a specific compiler
cp2p --input example.cpp --output ./bindings --compiler gcc

# Use a configuration file
cp2p --input example.cpp --output ./bindings --config config.json
```

### Command Line Arguments

- `--input`: Path to the C++ source file or project entry point
- `--output`: Output directory for generated bindings (default: ./bindings)
- `--compiler`: Compiler choice (gcc, clang, msvc, auto)
- `--config`: Optional JSON config file (if not provided, will parse C++ file)

### Configuration File Example

```json
{
  "classes": [
    {
      "name": "MyClass",
      "methods": ["method1", "method2"],
      "constructors": ["default", "withParams"]
    }
  ],
  "functions": ["globalFunc1", "globalFunc2"],
  "include_paths": ["/path/to/includes"],
  "compiler_flags": ["-std=c++17", "-O2"]
}
```

### C++ Code Example

```cpp
// example.cpp
#include <string>

class MyClass {
public:
    MyClass() {}
    std::string getMessage() { return "Hello from C++!"; }
};

// Bindings will be automatically generated for this class
```

### Using Generated Bindings

```python
# Python code using the generated bindings
from bindings import MyClass

obj = MyClass()
print(obj.getMessage())  # Output: Hello from C++!
```

## Compiler Detection

The tool automatically detects available compilers in the following order:

### Windows
1. MSVC (cl.exe)
2. GCC/MinGW (g++)
3. Clang (clang++)

### Linux/macOS
1. Clang (clang++)
2. GCC (g++)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test
go test ./compiler -run TestCompilerDetection
```

### Adding New Compiler Support

1. Add new compiler type to `CompilerType` enum
2. Implement detection logic in `detect.go`
3. Add compiler-specific flags and options
4. Update tests in `detect_test.go`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details
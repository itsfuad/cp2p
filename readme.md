# Cp2P: Seamless C++ to Python Binding Generator

Cp2P is a powerful CLI tool that automatically generates Python bindings for C++ code. It supports multiple compilers and platforms, making it easy to create Python interfaces for your C++ libraries.

## Features

- **Multi-Compiler Support**: Works with GCC, Clang, and MSVC
- **Cross-Platform**: Supports Windows, Linux, and macOS
- **Type-Safe**: Generates proper type hints and docstrings
- **Easy to Use**: Simple command-line interface
- **Configurable**: Supports JSON/YAML configuration files
- **Modern Python**: Generates Python 3.7+ compatible code

## Installation

### From Source

```bash
git clone https://github.com/itsfuad/Cp2P.git
cd Cp2P
go build
```

### Using Go

```bash
go install github.com/itsfuad/Cp2P@latest
```

## Quick Start

1. Create a C++ source file (`example.cpp`):

```cpp
extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
```

2. Generate Python bindings:

```bash
Cp2P --input example.cpp --output ./bindings
```

3. Use in Python:

```python
from bindings.example import add

result = add(5, 3)  # Returns 8
```

## Usage

### Basic Usage

```bash
Cp2P --input <source_file> --output <output_dir> [options]
```

### Options

- `--input`: Path to C++ source file (required)
- `--output`: Output directory (default: ./bindings)
- `--compiler`: Compiler choice (gcc, clang, msvc, auto)
- `--config`: JSON/YAML config file for custom bindings

### Configuration File Example

```json
{
    "functions": [
        {
            "name": "add",
            "return_type": "int",
            "parameters": [
                {"name": "a", "type": "int"},
                {"name": "b", "type": "int"}
            ],
            "docstring": "Adds two integers"
        }
    ]
}
```

## Development

### Prerequisites

- Go 1.21 or later
- C++ compiler (GCC, Clang, or MSVC)
- Python 3.7 or later

### Building

```bash
# Clone the repository
git clone https://github.com/itsfuad/Cp2P.git
cd Cp2P

# Install dependencies
go mod download

# Build
go build

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Project Structure

```
Cp2P/
├── cmd/
│   └── Cp2P.go           # CLI entry point
├── compiler/
│   ├── detect.go          # Compiler detection
│   └── compile.go         # Compilation logic
├── binding/
│   └── generator.go       # Python binding generator
├── config/
│   └── parser.go          # Config file parser
├── util/
│   ├── os_utils.go        # OS utilities
│   └── logger.go          # Logging
├── tests/                 # Test files
├── .github/              # GitHub Actions
├── README.md
├── go.mod
└── LICENSE
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by pybind11 and cppyy
- Built with Go and Python
- Uses ctypes for Python bindings


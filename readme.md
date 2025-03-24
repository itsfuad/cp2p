# Cpp2Py: Seamless C++ to Python Binding Generator

---

## Overview

Cpp2Py is a cross-platform CLI tool designed to:

- Compile C++ projects into shared libraries (DLL/SO/Dylib).
- Automatically generate Python bindings that feel native, with:
  - **Natural Imports**: Use functions as if they were pure Python.
  - **Intellisense Support**: Benefit from type hints and docstrings.
  - **No Manual Setup**: Ready-to-use bindings without additional configuration.

---

## Key Features

- **Multi-Compiler Support**: Compatible with GCC, Clang, and MSVC.
- **Cross-Platform Compatibility**: Supports Windows, Linux, and macOS.
- **Pythonic API**: Generated bindings are intuitive and integrate seamlessly with Python codebases.
- **Enhanced Developer Experience**: Auto-generated type hints and docstrings ensure full support for code editors and IDEs.

---

## Project Structure

```
cpp2py/
├── cmd/
│   └── cpp2py.go           # CLI entry point
├── compiler/
│   ├── detect.go           # Compiler detection
│   └── compile.go          # Cross-platform compilation logic
├── binding/
│   └── generator.go        # Python binding generator with type hints and docstrings
├── config/
│   └── parser.go           # Config file parser (JSON/YAML)
├── util/
│   ├── os_utils.go         # OS & path utilities
│   └── logger.go           # Logging & error handling
├── README.md
├── go.mod
└── LICENSE
```

---

## How It Works

### 1. Command Usage

```
cpp2py --input mylib.cpp --output ./bindings --compiler auto --config bindings.json
```

Arguments:

- `--input`: Path to the C++ source file or project entry point.
- `--output`: Output directory (default: `./bindings`).
- `--compiler`: Compiler choice (`gcc`, `clang`, `msvc`, `auto`).
- `--config`: (Optional) JSON/YAML config specifying functions to export.

---

### 2. Process Flow

1. **Compiler Detection**: Determines the available compiler based on the OS and `--compiler` flag.
2. **OS Detection**: Identifies the operating system to tailor the compilation process.
3. **Compilation**: Generates platform-specific shared libraries:
   - Windows: `.dll`
   - Linux: `.so`
   - macOS: `.dylib`
4. **Function Signature Handling**: Extracts function signatures and descriptions from:
   - Config files.
   - Annotated comments in C++ code:

```cpp
// EXPORT: int add(int a, int b) -> "Adds two integers."
```

5. **Python Binding Generation**: Produces a Python module with:
   - Natural function imports.
   - Type hints and docstrings for each function.
   - Encapsulated `ctypes` logic to maintain a clean API.
6. **Output Structure**:

```
/bindings/
├── libmylib.so     # Linux
├── libmylib.dylib  # macOS
├── mylib.dll       # Windows
└── mylib.py        # Python binding script
```

7. **Python Usage**: Developers can immediately utilize the bindings:

```python
from bindings.mylib import add

result = add(5, 7)
```

---

## Generated Python Binding Example

```python
import ctypes
import sys
import os
from typing import Any

# Load the shared library based on the OS
if sys.platform.startswith('win'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), 'mylib.dll'))
elif sys.platform.startswith('linux'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), 'libmylib.so'))
elif sys.platform.startswith('darwin'):
    _lib = ctypes.CDLL(os.path.join(os.path.dirname(__file__), 'libmylib.dylib'))

# Define the argument and return types for the C function
_lib.add.argtypes = [ctypes.c_int, ctypes.c_int]
_lib.add.restype = ctypes.c_int

def add(a: int, b: int) -> int:
    """
    Adds two integers.

    Args:
        a (int): First integer.
        b (int): Second integer.

    Returns:
        int: Sum of a and b.
    """
    return _lib.add(a, b)

__all__ = ['add']
```

---

## Benefits of This Approach

- **Pythonic Interface**: Users interact with functions as if they were native Python, without dealing with `ctypes` directly.
- **Enhanced Developer Experience**: Type hints and docstrings provide clarity and support for code editors, facilitating features like autocomplete and inline documentation.
- **Encapsulation**: The underlying `ctypes` implementation is hidden, offering a clean and intuitive API.

---

## Optional Features (Planned)

- **CMake Project Support**: Automatically detect and build CMake projects.
- **Signature Auto-Parsing**: Scan `.cpp` files for annotated comments to extract function signatures.
- **GUI Frontend**: Develop a user-friendly interface using Go’s Fyne library.
- **Virtualenv & Packaging**: Option to create Python packages for distribution.
- **Parallel Compilation**: Enhance performance for large C++ projects.

---

## Dependencies

- **Go Standard Library**: The core functionality is built using Go's standard packages.
- Optional: Utilize `cobra` for advanced CLI features.

---

## Project Timeline

*TBD*

---

## License

This project is licensed under the MIT License.

---

## Development Notes

- Ensure modular and clean code architecture.
- Maintain cross-platform consistency.
- Avoid external dependencies unless absolutely necessary.
- Provide comprehensive documentation and logging.
- Adhere to standard Go formatting practices (`gofmt`).

---

## Contribution Guidelines

1. Write modular and maintainable code.
2. Ensure cross-platform compatibility.
3. Minimize external dependencies.
4. Maintain clear documentation and logging practices.
5. Follow Go's standard formatting guidelines.

---

## Future Vision

Cpp2Py aims to be the definitive tool for bridging C++ and Python, allowing developers to integrate native code effortlessly, with bindings that feel as natural as writing pure Python.

---

**Would you like assistance in setting up the initial repository structure or any specific component of the project?**


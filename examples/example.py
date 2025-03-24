import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from bindings.math import add, multiply, square_root, is_even

def main():
    # Test add function
    result = add(5, 3)
    print(f"5 + 3 = {result}")

    # Test multiply function
    result = multiply(4.5, 2.0)
    print(f"4.5 * 2.0 = {result}")

    # Test square_root function
    result = square_root(16.0)
    print(f"sqrt(16.0) = {result}")

    # Test is_even function
    number = 42
    result = is_even(number)
    print(f"Is {number} even? {result}")

if __name__ == "__main__":
    main() 
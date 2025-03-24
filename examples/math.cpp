#include <cmath>

extern "C" {

// EXPORT: int add(int a, int b) -> "Adds two integers."
int add(int a, int b) {
    return a + b;
}

// EXPORT: double multiply(double a, double b) -> "Multiplies two floating-point numbers."
double multiply(double a, double b) {
    return a * b;
}

// EXPORT: double square_root(double x) -> "Calculates the square root of a number."
double square_root(double x) {
    return std::sqrt(x);
}

// EXPORT: bool is_even(int n) -> "Checks if a number is even."
bool is_even(int n) {
    return n % 2 == 0;
}

} // extern "C" 
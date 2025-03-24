package util

import (
	"fmt"
	"log"
	"os"
)

// Logger provides logging functionality for the application
type Logger struct {
	*log.Logger
	debug bool
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		debug:  false,
	}
}

// SetDebug enables or disables debug logging
func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		l.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.Printf("[WARN] "+format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

// Fatalf logs a fatal message with formatting and exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
}

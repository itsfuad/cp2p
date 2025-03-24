package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the binding configuration
type Config struct {
	Functions []FunctionConfig `json:"functions"`
	Includes  []string         `json:"includes"`
	Libraries []string         `json:"libraries"`
	Types     []TypeConfig     `json:"types"` // Complex types (structs, classes, etc.)
}

// TypeConfig represents a complex type definition
type TypeConfig struct {
	Name        string   `json:"name"`        // Name of the type
	Kind        string   `json:"kind"`        // struct, class, enum, union
	Fields      []Field  `json:"fields"`      // For structs/classes
	Values      []string `json:"values"`      // For enums
	BaseType    string   `json:"base_type"`   // For enums
	Description string   `json:"description"` // Documentation
}

// Field represents a field in a struct/class
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// FunctionConfig represents the configuration for a single function
type FunctionConfig struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Parameters  []Param `json:"parameters"`
	ReturnType  string  `json:"return_type"`
	Docstring   string  `json:"docstring"`
}

// Param represents a function parameter
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ParseConfig parses a JSON configuration file
func ParseConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %v", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if len(cfg.Functions) == 0 {
		return fmt.Errorf("no functions specified in config")
	}

	for i, fn := range cfg.Functions {
		if fn.Name == "" {
			return fmt.Errorf("function at index %d has no name", i)
		}
		if fn.ReturnType == "" {
			return fmt.Errorf("function %s has no return type", fn.Name)
		}
	}

	return nil
}

// GetFunctionConfig returns the configuration for a specific function
func (c *Config) GetFunctionConfig(name string) *FunctionConfig {
	for _, fn := range c.Functions {
		if fn.Name == name {
			return &fn
		}
	}
	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Functions: []FunctionConfig{},
		Includes:  []string{},
		Libraries: []string{},
	}
}

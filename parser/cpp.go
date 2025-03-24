package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"cp2p/config"
)

// ParseCppFile parses a C++ file and extracts functions marked with EXPORT comments
func ParseCppFile(filePath string) (*config.Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var functions []config.FunctionConfig
	exportRegex := regexp.MustCompile(`//\s*EXPORT:\s*(\w+)\s+(\w+)\s*\((.*?)\)\s*->\s*"([^"]*)"`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := exportRegex.FindStringSubmatch(line)
		if matches != nil {
			// matches[1] = return type
			// matches[2] = function name
			// matches[3] = parameters
			// matches[4] = description
			fn := config.FunctionConfig{
				Name:        matches[2],
				Description: matches[4],
				ReturnType:  matches[1],
				Parameters:  parseParameters(matches[3]),
			}
			functions = append(functions, fn)
		}
	}

	return &config.Config{
		Functions: functions,
		Includes:  []string{},
		Libraries: []string{},
	}, nil
}

func parseParameters(paramStr string) []config.Param {
	if paramStr == "" {
		return []config.Param{}
	}

	params := strings.Split(paramStr, ",")
	var result []config.Param

	for _, p := range params {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Split type and name
		parts := strings.Fields(p)
		if len(parts) >= 2 {
			paramType := parts[0]
			paramName := parts[1]
			// Remove any trailing semicolons or other characters
			paramName = strings.TrimRight(paramName, ";")

			result = append(result, config.Param{
				Name:        paramName,
				Type:        paramType,
				Description: "", // Could be enhanced to parse parameter descriptions from comments
			})
		}
	}

	return result
}

// Package goenvify provides a simple way to load environment variables
package goenvify

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const EnvPartsNumber = 2

var ErrEmptyFile = errors.New("file is empty")

// LoadFile loads the content of a file and returns it as a string.
// It returns an error if the file does not exist or is empty.
func LoadFile(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("could not read file %s: %w", file, err)
	}

	if len(content) == 0 {
		return "", ErrEmptyFile
	}

	return string(content), nil
}

// SplitContent splits the content of a file by line and returns a slice of strings.
// It removes comments and empty lines.
func SplitContent(content string) []string {
	var lines []string

	for line := range strings.SplitSeq(content, "\n") {
		currentLine := line
		trimComment := strings.Split(currentLine, " #")

		if len(trimComment) > 1 {
			currentLine = trimComment[0]
		}

		currentLine = strings.TrimSpace(currentLine)
		if currentLine != "" && !strings.HasPrefix(currentLine, "#") {
			lines = append(lines, currentLine)
		}
	}

	return lines
}

// MapContent maps the content of a file to a map of strings.
// It removes empty strings and trims spaces.
func MapContent(content []string) map[string]string {
	envVars := make(map[string]string)
	sepRegex := regexp.MustCompile(`[=:]`)
	quoteRegex := regexp.MustCompile(`['"]`)

	for _, line := range content {
		if line == "" {
			continue
		}

		if !sepRegex.MatchString(line) {
			continue
		}

		envVar := sepRegex.Split(line, EnvPartsNumber)
		envKey := strings.TrimSpace(envVar[0])
		envVal := strings.TrimSpace(envVar[1])

		if envKey != "" {
			envVars[envKey] = quoteRegex.ReplaceAllString(envVal, "")
		}
	}

	return envVars
}

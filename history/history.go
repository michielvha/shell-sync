package history

import (
	"bufio"
	"os"
)

// ReadLines reads all lines from a file.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines writes all lines to a file (atomic write).
func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// MergeHistories merges two history lists, preserving order as much as possible and deduplicating.
func MergeHistories(local, remote []string) []string {
	seen := make(map[string]struct{})
	merged := []string{}
	for _, line := range local {
		if _, exists := seen[line]; !exists {
			merged = append(merged, line)
			seen[line] = struct{}{}
		}
	}
	for _, line := range remote {
		if _, exists := seen[line]; !exists {
			merged = append(merged, line)
			seen[line] = struct{}{}
		}
	}
	return merged
}

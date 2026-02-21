package functions

import (
	"fmt"
	"os"
)

// ReadFile reads the contents of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(data), nil
}

func CreateNewFile(filename string) string {
	err := os.WriteFile(filename, nil, 0644)
	if err != nil {
		fmt.Errorf("error writing file: %w", err)
	}

	return filename
}

func UpdateFile(filename string, content string) (string, error) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return "Saved successfully", nil
}

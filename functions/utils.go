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

// WriteFile writes content to a file
func writeFile(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func CreateNewFile(filename string) string {
	err := os.WriteFile(filename, nil, 0644)
	if err != nil {
		fmt.Errorf("error writing file: %w", err)
	}
	
	return filename
}

func ShowFiles(path string) ([]string, error) {
	var dirArr []string
	data, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, pth := range data {
		if pth.IsDir() {
			dirArr = append(dirArr, pth.Name())
		}
	}
	return dirArr, nil
}

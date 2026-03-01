package functions

import (
	"fmt"
	"os"
)

// TODO UPDATE ("/Users/user/Desktop/Appunti/"+filename in CWD
// ReadFile reads the contents of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(data), nil
}

func CreateNewFile(filename string) string {
	//TODO fix empty name area
	//TODO auto create folder
	err := os.WriteFile("C:\\Users\\Francesco\\Desktop\\Appunti\\"+filename, nil, 0644)
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

func CompileMarkdown(currentFile string) {

}

package utils

import "os"

func CreateFileIfNotExists(filePath string) error {
	// Check if the file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return nil // File already exists, no need to create it
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err // Return error if file creation fails
	}
	defer file.Close() // Ensure the file is closed after creation

	return nil // Successfully created the file
}

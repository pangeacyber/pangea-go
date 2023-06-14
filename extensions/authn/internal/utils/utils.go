package utils

import (
	"fmt"
	"os"
)

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("file %v does not exist\n", filePath)
		return false
	}

	return true
}

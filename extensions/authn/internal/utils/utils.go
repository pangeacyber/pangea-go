package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("file %v does not exist\n", filePath)
		return false
	}
	return true
}

func ValidateAndOpen(fileName string) (*os.File, error) {
	if !IsFileExist(fileName) {
		return nil, errors.New("mapping file does not exist")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ConvertValue(value interface{}, sourceType string, destType string) (interface{}, error) {
	// If source field is not defined then return as it is
	if sourceType == "" {
		return value, nil
	}
	// conversion
	// TODO - Add more conversion
	switch sourceType {
	case "string":
		switch destType {
		case "integer":
			return strconv.Atoi(value.(string))
		default:
			return value, nil
		}
	case "integer":
		switch destType {
		case "string":
			return strconv.Itoa(value.(int)), nil
		default:
			return value, nil
		}
	default:
		return nil, fmt.Errorf("unsupported type conversion: %s to %s", sourceType, destType)
	}
}

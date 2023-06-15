package models

import (
	"encoding/json"
	"errors"
	"extensions/authn/internal/utils"
	"fmt"
	"os"
)

// Mapping represents the structure of our mapping file.
type MappedField struct {
	SourceField      string `json:"sourceField"`
	DestField        string `json:"destField"`
	SourceType       string `json:"sourceType,omitempty"`
	DestType         string `json:"destType,omitempty"`
	DestValueExtract string `json:"destValueExtract,omitempty""`
}

type Mappings struct {
	fields []MappedField
}

func NewMappings(mappingFile string) (*Mappings, error) {
	if !utils.IsFileExist(mappingFile) {
		return nil, errors.New(fmt.Sprintf("file=%s does not exist", mappingFile))
	}
	fileData, err := os.ReadFile(mappingFile)
	if err != nil {
		// Handle the error
		return nil, err
	}
	// Unmarshal the mapping file into a slice of Mappings.
	var mappings []MappedField
	if err := json.Unmarshal(fileData, &mappings); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &Mappings{
		fields: mappings,
	}, nil
}

func (mapping Mappings) MappedFields(data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(data))
	for _, mapping := range mapping.fields {
		value, exists := data[mapping.SourceField]
		if !exists {
			fmt.Printf("field %s did not map, copying as it is", mapping.SourceType)
			result[mapping.SourceField] = value
			continue
		}

		convertedValue, err := utils.ConvertValue(value, mapping.SourceType, mapping.DestType)
		if err != nil {
			fmt.Printf("Could not convert field %s: %v, skipping it\n", mapping.SourceField, err)
			continue
		}
		result[mapping.DestField] = convertedValue
	}
	return result, nil
}

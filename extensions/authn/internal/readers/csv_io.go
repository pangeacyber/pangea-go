package readers

import (
	"encoding/csv"
	"errors"
	"extensions/authn/internal/models"
	"extensions/authn/internal/utils"
	"fmt"
	"os"
)

// CSVReader A CSV file reader
type CSVReader struct {
	csvReader *csv.Reader // CSV file path
	mapping   *models.Mappings
	headers   []string
}

// NewCSVReader
// Assumption that provided csv file has header row
// Create new CSV readers
// filePath - full csv file path
//
// */
func NewCSVReader(filePath string, mapping *models.Mappings) (*CSVReader, error) {
	if filePath == "" {
		return nil, errors.New("empty file is not accepted")
	}
	if !utils.IsFileExist(filePath) {
		return nil, fmt.Errorf("file %v does not exist", filePath)
	}
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("failed to open file, err=%v", err)
		return nil, err
	}
	csvReader := csv.NewReader(csvFile)
	var headerRow []string
	headerRow, err = csvReader.Read()
	if err != nil {
		fmt.Errorf("failed to read header")
		return nil, err
	}

	return &CSVReader{
		csvReader: csvReader,
		headers:   headerRow,
		mapping:   mapping,
	}, nil
}

// Next
// @summary Iterator to read single csv record
// @return single record with headers
// */
func (csvReader *CSVReader) Next() (map[string]interface{}, error) {
	record, err := csvReader.csvReader.Read()
	if err != nil {
		return nil, err
	}
	// Construct map using headers
	rowData := make(map[string]interface{}, len(record))
	for idx := range record {
		rowData[csvReader.headers[idx]] = record[idx]
	}
	fmt.Println(rowData)
	if csvReader.mapping != nil {
		rowData, err = csvReader.mapping.MappedFields(rowData)
		if err != nil {
			return nil, err
		}
	}
	return rowData, nil
}

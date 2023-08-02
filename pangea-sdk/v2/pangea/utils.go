package pangea

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func HashSHA256(i string) string {
	b := sha256.Sum256([]byte(i))
	return hex.EncodeToString(b[:])
}

func GetHashPrefix(h string, len uint) string {
	return h[0:len]
}

func initFileWriter() {
	// Open the output file
	filename := "pangea_sdk_log.json"
	var err error
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// Where should we close this file?
	if err != nil {
		fmt.Printf("Failed to open log file: %s. Logger will go to stdout", filename)
		file = os.Stdout
	}
}

func GetDefaultPangeaLogger() *zerolog.Logger {
	// Set up the logger
	initFileWriterOnce.Do(initFileWriter)

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000Z07:00"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"

	// Set up the JSON file writer as the output
	logger := zerolog.New(file).With().Timestamp().Logger()
	return &logger
}

package pangea

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
)

func HashSHA256(i string) string {
	b := sha256.Sum256([]byte(i))
	return hex.EncodeToString(b[:])
}

func HashSHA1(i string) string {
	b := sha1.Sum([]byte(i))
	return hex.EncodeToString(b[:])
}

func HashSHA512(i string) string {
	b := sha512.Sum512([]byte(i))
	return hex.EncodeToString(b[:])
}

func GetHashPrefix(h string, len uint) string {
	return h[0:len]
}

type UploadFileParams struct {
	Size   int    `json:"size,omitempty"`
	CRC    string `json:"crc32c,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

func GetUploadFileParams(file *os.File) (*UploadFileParams, error) {
	// Create a new CRC32C hash
	crcHash := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	// Create a new SHA256 hash
	sha256Hash := sha256.New()

	// Seek back to the beginning of the file
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Copy the file content into the hash function
	if _, err := io.Copy(sha256Hash, file); err != nil {
		return nil, err
	}
	// Get the hash sum as a byte slice
	hashInBytes := sha256Hash.Sum(nil)

	// Seek back to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Copy the file content into the hash calculation
	size, err := io.Copy(crcHash, file)
	if err != nil {
		return nil, err
	}
	// Get the CRC32C checksum value
	crc32c := crcHash.Sum32()

	// Reset to be sent
	file.Seek(0, 0)

	// Convert the CRC32 value to hexadecimal
	crcStr := strconv.FormatUint(uint64(crc32c), 16)
	// Pad "0" on the left to make it 8 characters long. It's for the zero bytes file case
	paddedCRCStr := strings.Repeat("0", 8-len(crcStr)) + crcStr

	return &UploadFileParams{
		CRC:    paddedCRCStr,
		SHA256: hex.EncodeToString(hashInBytes),
		Size:   int(size),
	}, nil
}

func GetFileSize(file *os.File) (int64, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	// Seek to the end of the file
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Reset to be sent
	file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, nil
}

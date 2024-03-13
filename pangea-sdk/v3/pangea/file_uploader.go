package pangea

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"strconv"
)

type FileUploader struct {
	client *Client
}

func NewFileUploader() FileUploader {
	cfg := &Config{
		QueuedRetryEnabled: false,
	}

	return FileUploader{
		client: NewClient("FileUploader", cfg),
	}
}

func (fu *FileUploader) UploadFile(ctx context.Context, url string, tm TransferMethod, fd FileData) error {
	if tm == TMmultipart {
		return fmt.Errorf("%s is not supported in UploadFile. Use service client instead", tm)
	}

	fds := FileData{
		File:    fd.File,
		Name:    "file",
		Details: fd.Details,
	}
	return fu.client.UploadFile(ctx, url, tm, fds)
}

type UploadFileParams struct {
	Size   int    `json:"size,omitempty"`
	CRC32C string `json:"crc32c,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

func GetUploadFileParams(input io.ReadSeeker) (*UploadFileParams, error) {
	// Create a new CRC32C hash
	crcHash := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	// Create a new SHA256 hash
	sha256Hash := sha256.New()

	// Seek back to the beginning of the file
	_, err := input.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Copy the file content into the hash function
	if _, err := io.Copy(sha256Hash, input); err != nil {
		return nil, err
	}
	// Get the hash sum as a byte slice
	hashInBytes := sha256Hash.Sum(nil)

	// Seek back to the beginning of the file
	_, err = input.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Copy the file content into the hash calculation
	size, err := io.Copy(crcHash, input)
	if err != nil {
		return nil, err
	}
	// Get the CRC32C checksum value
	crc32c := crcHash.Sum32()

	// Reset to be sent
	input.Seek(0, 0)

	return &UploadFileParams{
		CRC32C: strconv.FormatUint(uint64(crc32c), 16),
		SHA256: hex.EncodeToString(hashInBytes),
		Size:   int(size),
	}, nil
}

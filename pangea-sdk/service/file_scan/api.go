package file_scan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strconv"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

// @summary Scan
//
// @description Scan a file for malicious content.
//
// @operationId file_scan_post_v1_scan
//
// @example
//
//	input := &file_scan.FileScanRequest{
//		Raw: true,
//		Verbose: true,
//		Provider: "crowdstrike",
//	}
//
//	// This should be your own file to scan
//	file, err := os.Open("./path/to/file.pdf")
//	if err != nil {
//		log.Fatal("expected no error got: %v", err)
//	}
//
//	resp, err := client.Scan(ctx, input, file)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
func (e *FileScan) Scan(ctx context.Context, input *FileScanRequest, file *os.File) (*pangea.PangeaResponse[FileScanResult], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	var req FileScanFullRequest
	params := &FileScanFileParams{}

	if input.TransferMethod == pangea.TMpostURL {
		var err error
		params, err = GetUploadFileParams(file)
		if err != nil {
			return nil, err
		}
	}

	req = FileScanFullRequest{
		FileScanRequest:    *input,
		FileScanFileParams: *params,
	}

	name := "file"
	if input.TransferMethod == pangea.TMmultipart {
		name = "upload"
	}

	fd := pangea.FileData{
		File: file,
		Name: name,
	}

	return request.DoPostWithFile(ctx, e.Client, "v1/scan", &req, &FileScanResult{}, fd)
}

func (e *FileScan) RequestUploadURL(ctx context.Context, input *FileScanGetURLRequest, file *os.File) (*pangea.PangeaResponse[FileScanResult], error) {
	if input.TransferMethod == pangea.TMpostURL && input.FileParams == nil {
		return nil, errors.New("Need to set FileParams in order to use TMpostURL or TMdirect")
	}

	req := &FileScanFullRequest{
		FileScanRequest: FileScanRequest{
			TransferRequest: pangea.TransferRequest{
				TransferMethod: input.TransferMethod,
			},
		},
	}
	if input.FileParams != nil {
		req.FileScanFileParams = *input.FileParams
	}

	return request.GetUploadURL(ctx, e.Client, "v1/scan", req, &FileScanResult{})
}

type FileScanRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest
	pangea.TransferRequest

	Verbose   bool   `json:"verbose,omitempty"`
	Raw       bool   `json:"raw,omitempty"`
	Provider  string `json:"provider,omitempty"`
	SourceURL string `json:"source_url,omitempty"`
}

type FileScanFileParams struct {
	Size   int    `json:"size,omitempty"`
	CRC    string `json:"crc32c,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

type FileScanGetURLRequest struct {
	TransferMethod pangea.TransferMethod
	Verbose        bool
	Raw            bool
	Provider       string
	FileParams     *FileScanFileParams
}

type FileScanFullRequest struct {
	FileScanRequest
	FileScanFileParams
}

type FileScanData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type FileScanResult struct {
	Data       FileScanData `json:"data"`
	Parameters interface{}  `json:"parameters,omitempty"`
	RawData    interface{}  `json:"raw_data,omitempty"`
}

func GetUploadFileParams(file *os.File) (*FileScanFileParams, error) {
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

	return &FileScanFileParams{
		CRC:    strconv.FormatUint(uint64(crc32c), 16),
		SHA256: hex.EncodeToString(hashInBytes),
		Size:   int(size),
	}, nil

}

func (fu *FileUploader) UploadFile(ctx context.Context, url string, tm pangea.TransferMethod, fd pangea.FileData) error {
	if tm == pangea.TMmultipart {
		return errors.New(fmt.Sprintf("%s is not supported in UploadFile. Use Scan() instead", tm))
	}

	fds := pangea.FileData{
		File:    fd.File,
		Name:    "file",
		Details: fd.Details,
	}
	return fu.client.UploadFile(ctx, url, tm, fds)
}

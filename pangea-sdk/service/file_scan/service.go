package file_scan

import (
	"context"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	Scan(ctx context.Context, input *FileScanRequest, file *os.File) (*pangea.PangeaResponse[FileScanResult], error)
	RequestUploadURL(ctx context.Context, input *FileScanGetURLRequest, file *os.File) (*pangea.PangeaResponse[FileScanResult], error)

	// Base service methods
	pangea.BaseServicer
}

type FileScan struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &FileScan{
		BaseService: pangea.NewBaseService("file-scan", cfg),
	}
	return cli
}

type FileUploader struct {
	client *pangea.Client
}

func NewFileUploader() FileUploader {
	cfg, err := pangea.NewConfig(option.WithQueuedRetryEnabled(false))
	if err != nil {
		panic(err)
	}

	return FileUploader{
		client: pangea.NewClient("FileScanUploader", cfg),
	}
}

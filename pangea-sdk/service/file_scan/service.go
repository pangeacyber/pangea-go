package file_scan

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Scan(ctx context.Context, input *FileScanRequest, file io.Reader) (*pangea.PangeaResponse[FileScanResult], error)
}

type FileScan struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &FileScan{
		Client: pangea.NewClient("file-scan", cfg),
	}
	return cli
}
